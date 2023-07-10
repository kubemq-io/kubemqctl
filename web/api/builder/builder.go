package builder

import (
	"context"
	"fmt"

	"strconv"
)

// const manifestLinkFormat = "http://localhost:10100/build/%d"
const kubectlPrompt = "kubectl apply -f "
const manifestLinkFormat = "https://deploy.kubemq.io/build/%d"
const initLinkFormat = "https://deploy.kubemq.io/init"

type Service struct {
	repository *Repository
}

func NewBuilderService(ctx context.Context) *Service {
	return &Service{
		repository: NewRepository(ctx),
	}
}

func (s *Service) ProcessRequest(req *BuilderRequest) (*BuilderResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	resp := NewBuilderResponse()
	if req.Model.SetInit {
		resp.AddLink("", kubectlPrompt+initLinkFormat, initLinkFormat)
	}
	mb := NewManifest().SetID(req.Id).SetType(req.Type)
	switch req.Type {
	case "bridges", "sources", "targets":
		connector := req.Connector()
		if err := connector.Validate(); err != nil {
			return nil, err
		}
		if req.Model.SetOperator {
			om, _ := NewOperatorManifest().SetNamespace(connector.Metadata.Namespace).Manifest()
			mb.AddPrepend(om)
		}
		mb.SetItems(connector)
	case "clusters":
		clusters, namespaces := req.Clusters()
		if req.Model.SetOperator {
			for _, ns := range namespaces {
				om, _ := NewOperatorManifest().SetNamespace(ns).Manifest()
				mb.AddPrepend(om)
			}
		}
		for _, cluster := range clusters {
			mb.SetItems(cluster)
		}
		if len(clusters) > 0 {
			resp.SetKey(clusters[0].Key())
		}
	default:
		return nil, fmt.Errorf("invalid request type")
	}
	s.repository.Set(mb)
	link := fmt.Sprintf(manifestLinkFormat, mb.Id)
	resp.AddLink("", kubectlPrompt+link, link)
	return resp, nil
}

func (s *Service) GetBuild(id string) (string, error) {
	n, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return "", fmt.Errorf("invalid build id")
	}
	m, err := s.repository.Get(n)
	if err != nil {
		return "", fmt.Errorf("build not found")
	}
	return m.String(), nil

}
