package services_test

import (
	"errors"
	"io/ioutil"
	"os"

	"gp_upgrade/helpers"
	"gp_upgrade/hub/configutils"
	"gp_upgrade/hub/services"
	pb "gp_upgrade/idl"
	"gp_upgrade/utils"

	"google.golang.org/grpc"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("PrepareShutdownClusters", func() {
	// ignoring the go routine
	It("returns successfully", func() {
		utils.System.RemoveAll = func(s string) error { return nil }
		utils.System.MkdirAll = func(s string, perm os.FileMode) error { return nil }

		reader := configutils.NewReader()
		dir, err := ioutil.TempDir("", "")
		Expect(err).ToNot(HaveOccurred())
		conf := &services.HubConfig{
			StateDir: dir,
		}
		hub := services.NewHub(&mockClusterPair{}, &reader, grpc.DialContext, nil, conf)

		_, err = hub.PrepareShutdownClusters(nil, &pb.PrepareShutdownClustersRequest{})
		Expect(err).To(BeNil())
	})

	It("fails if the cluster configuration setup can't be loaded", func() {
		utils.System.RemoveAll = func(s string) error { return nil }
		utils.System.MkdirAll = func(s string, perm os.FileMode) error { return nil }

		reader := configutils.NewReader()
		dir, err := ioutil.TempDir("", "")
		Expect(err).ToNot(HaveOccurred())
		conf := &services.HubConfig{
			StateDir: dir,
		}
		clusterPair := &mockClusterPair{
			InitErr: errors.New("boom"),
		}
		hub := services.NewHub(clusterPair, &reader, grpc.DialContext, nil, conf)

		_, err = hub.PrepareShutdownClusters(nil, &pb.PrepareShutdownClustersRequest{})
		Expect(err).To(MatchError("boom"))
	})
})

type mockClusterPair struct {
	InitErr error
}

func (c *mockClusterPair) StopEverything(str string) {}
func (c *mockClusterPair) Init(baseDir, oldPath, newPath string, execer helpers.CommandExecer) error {
	return c.InitErr
}
func (c *mockClusterPair) GetPortsAndDataDirForReconfiguration() (int, int, string) { return -1, -1, "" }
