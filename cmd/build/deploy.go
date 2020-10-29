package build

//
//import (
//	"fmt"
//	"github.com/kubemq-hub/builder/survey"
//	"github.com/kubemq-io/kubemqctl/pkg/k8s/client"
//	"github.com/kubemq-io/kubemqctl/pkg/k8s/manager/cluster"
//	"github.com/kubemq-io/kubemqctl/pkg/k8s/manager/connector"
//	"github.com/kubemq-io/kubemqctl/pkg/k8s/manager/operator"
//	operatorTypes "github.com/kubemq-io/kubemqctl/pkg/k8s/types/operator"
//	"github.com/kubemq-io/kubemqctl/pkg/utils"
//)
//
//type deploy struct {
//	client     *client.Client
//	clusters   *ClustersBuilder
//	connectors *ConnectorsBuilder
//}
//
//func newDeploy() *deploy {
//	return &deploy{}
//}
//
//func (d *deploy) SetClient(value *client.Client) *deploy {
//	d.client = value
//	return d
//}
//func (d *deploy) SetClusters(value *ClustersBuilder) *deploy {
//	d.clusters = value
//	return d
//}
//func (d *deploy) SetConnectors(value *ConnectorsBuilder) *deploy {
//	d.connectors = value
//	return d
//}
//func (d *deploy) deployOperators() error {
//	namespaces := make(map[string]string)
//	for _, dep := range d.clusters.deployments {
//		namespaces[dep.Namespace] = dep.Namespace
//	}
//	for _, dep := range d.connectors.deployments {
//		namespaces[dep.Namespace] = dep.Namespace
//	}
//	operatorManager, err := operator.NewManager(d.client)
//	if err != nil {
//		return err
//	}
//	for _, ns := range namespaces {
//		if !operatorManager.IsKubemqOperatorExists(ns) {
//			operatorDeployment, err := operatorTypes.CreateDeployment("kubemq-operator", ns)
//			if err != nil {
//				return err
//			}
//			_, _, err = operatorManager.CreateOrUpdateKubemqOperator(operatorDeployment)
//			if err != nil {
//				return err
//			}
//			utils.Printlnf("Kubemq operator %s/kubemq-operator created.", ns)
//		} else {
//			utils.Printlnf("Kubemq operator %s/kubemq-operator exists", ns)
//		}
//	}
//	return nil
//}
//func (d *deploy) deployClusters() error {
//	numClusters := len(d.clusters.deployments)
//	if numClusters == 0 {
//		return nil
//	}
//	utils.Printlnf("Deploying %d KubeMQ Clusters...", numClusters)
//	clusterManager, err := cluster.NewManager(d.client)
//	if err != nil {
//		return err
//	}
//	for _, dep := range d.clusters.deployments {
//		cluster, isUpdate, err := clusterManager.CreateOrUpdateKubemqCluster(dep)
//		if err != nil {
//			return err
//		}
//		if isUpdate {
//			utils.Printlnf("kubemq cluster %s/%s configured.", cluster.Namespace, cluster.Name)
//		} else {
//			utils.Printlnf("kubemq cluster %s/%s created.", cluster.Namespace, cluster.Name)
//		}
//	}
//	return nil
//}
//func (d *deploy) deployConnectors() error {
//	numConnectors := len(d.connectors.deployments)
//	if numConnectors == 0 {
//		return nil
//	}
//	utils.Printlnf("Deploying %d KubeMQ Connectors...", numConnectors)
//	connectorManager, err := connector.NewManager(d.client)
//	if err != nil {
//		return err
//	}
//	for _, dep := range d.connectors.deployments {
//		connector, isUpdate, err := connectorManager.CreateOrUpdateKubemqConnector(dep)
//		if err != nil {
//			return err
//		}
//		if isUpdate {
//			utils.Printlnf("kubemq connector %s/%s configured.", connector.Namespace, connector.Name)
//		} else {
//			utils.Printlnf("kubemq connector %s/%s created.", connector.Namespace, connector.Name)
//		}
//	}
//	return nil
//}
//func (d *deploy) confirmDeploy() (bool, error) {
//
//	numDeployments := len(d.connectors.deployments) + len(d.clusters.deployments)
//	if numDeployments == 0 {
//		return false, fmt.Errorf("no KubeMQ components found to deploy")
//	}
//	confirm := false
//	err := survey.NewBool().
//		SetKind("bool").
//		SetName("confirm").
//		SetMessage(fmt.Sprintf("Deploying %d KubeMQ Cluster and %d Connectors, are you sure", len(d.clusters.deployments), len(d.connectors.deployments))).
//		SetDefault("true").
//		SetRequired(true).
//		Render(&confirm)
//	return confirm, err
//
//}
//
//func (d *deploy) Do() error {
//
//	if confirm, err := d.confirmDeploy(); err != nil {
//		return err
//	} else {
//		if confirm {
//			if err := d.deployOperators(); err != nil {
//				return err
//			}
//			if err := d.deployClusters(); err != nil {
//				return err
//			}
//			if err := d.deployConnectors(); err != nil {
//				return err
//			}
//		}
//	}
//
//	return nil
//}
