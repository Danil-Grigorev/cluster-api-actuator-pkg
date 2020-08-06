package operators

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/openshift/cluster-api-actuator-pkg/pkg/framework"
	"k8s.io/utils/pointer"
)

var (
	maoDeployment        = "machine-api-operator"
	maoManagedDeployment = "machine-api-controllers"
	maoManagedDaemonSet  = "machine-api-termination-handler"
)

var _ = Describe("[Feature:Operators] Machine API operator deployment should", func() {
	defer GinkgoRecover()

	BeforeEach(func() {
		client, err := framework.LoadClient()
		Expect(err).NotTo(HaveOccurred())
		Expect(framework.IsStatusAvailable(client, "machine-api")).To(BeTrue())
	})

	It("be available", func() {
		client, err := framework.LoadClient()
		Expect(err).NotTo(HaveOccurred())
		Expect(framework.IsDeploymentAvailable(client, maoDeployment, framework.MachineAPINamespace)).To(BeTrue())
	})

	It("reconcile controllers deployment", func() {
		client, err := framework.LoadClient()
		Expect(err).NotTo(HaveOccurred())

		initialDeployment, err := framework.GetDeployment(client, maoManagedDeployment, framework.MachineAPINamespace)
		Expect(err).NotTo(HaveOccurred())

		By(fmt.Sprintf("checking deployment %q is available", maoManagedDeployment))
		Expect(framework.IsDeploymentAvailable(client, maoManagedDeployment, framework.MachineAPINamespace)).To(BeTrue())

		By(fmt.Sprintf("deleting deployment %q", maoManagedDeployment))
		Expect(framework.DeleteDeployment(client, initialDeployment)).NotTo(HaveOccurred())

		By(fmt.Sprintf("checking deployment %q is available again", maoManagedDeployment))
		Expect(framework.IsDeploymentAvailable(client, maoManagedDeployment, framework.MachineAPINamespace)).To(BeTrue())

		By(fmt.Sprintf("checking deployment %q spec matches", maoManagedDeployment))
		Expect(framework.IsDeploymentSynced(client, initialDeployment, maoManagedDeployment, framework.MachineAPINamespace)).To(BeTrue())
	})

	It("maintains deployment spec", func() {
		client, err := framework.LoadClient()
		Expect(err).NotTo(HaveOccurred())

		initialDeployment, err := framework.GetDeployment(client, maoManagedDeployment, framework.MachineAPINamespace)
		Expect(err).NotTo(HaveOccurred())

		By(fmt.Sprintf("checking deployment %q is available", maoManagedDeployment))
		Expect(framework.IsDeploymentAvailable(client, maoManagedDeployment, framework.MachineAPINamespace)).To(BeTrue())

		changedDeployment := initialDeployment.DeepCopy()
		changedDeployment.Spec.Replicas = pointer.Int32Ptr(0)

		By(fmt.Sprintf("updating deployment %q", maoManagedDeployment))
		Expect(framework.UpdateDeployment(client, maoManagedDeployment, framework.MachineAPINamespace, changedDeployment)).NotTo(HaveOccurred())

		By(fmt.Sprintf("checking deployment %q spec matches", maoManagedDeployment))
		Expect(framework.IsDeploymentSynced(client, initialDeployment, maoManagedDeployment, framework.MachineAPINamespace)).To(BeTrue())

		By(fmt.Sprintf("checking deployment %q is available again", maoManagedDeployment))
		Expect(framework.IsDeploymentAvailable(client, maoManagedDeployment, framework.MachineAPINamespace)).To(BeTrue())

	})

	It("reconcile termination handler daemonSet", func() {
		client, err := framework.LoadClient()
		Expect(err).NotTo(HaveOccurred())

		initialDaemonSet, err := framework.GetDaemonSet(client, maoManagedDaemonSet, framework.MachineAPINamespace)
		Expect(err).NotTo(HaveOccurred())

		By(fmt.Sprintf("checking daemonSet %q is available", maoManagedDaemonSet))
		Expect(framework.IsDaemonSetAvailable(client, maoManagedDaemonSet, framework.MachineAPINamespace)).To(BeTrue())

		By(fmt.Sprintf("deleting daemonSet %q", maoManagedDaemonSet))
		Expect(framework.DeleteDaemonSet(client, initialDaemonSet)).NotTo(HaveOccurred())

		By(fmt.Sprintf("checking daemonSet %q is available again", maoManagedDaemonSet))
		Expect(framework.IsDaemonSetAvailable(client, maoManagedDaemonSet, framework.MachineAPINamespace)).To(BeTrue())

		By(fmt.Sprintf("checking daemonSet %q spec matches", maoManagedDaemonSet))
		Expect(framework.IsDaemonSetSynced(client, initialDaemonSet, maoManagedDaemonSet, framework.MachineAPINamespace)).To(BeTrue())
	})

	It("maintains termination handler daemonSet spec", func() {
		client, err := framework.LoadClient()
		Expect(err).NotTo(HaveOccurred())

		initialDaemonSet, err := framework.GetDaemonSet(client, maoManagedDaemonSet, framework.MachineAPINamespace)
		Expect(err).NotTo(HaveOccurred())

		By(fmt.Sprintf("checking daemonSet %q is available", maoManagedDaemonSet))
		Expect(framework.IsDaemonSetAvailable(client, maoManagedDaemonSet, framework.MachineAPINamespace)).To(BeTrue())

		changedDaemonSet := initialDaemonSet.DeepCopy()
		changedDaemonSet.Spec.Selector = nil

		By(fmt.Sprintf("updating daemonSet %q", maoManagedDaemonSet))
		Expect(framework.UpdateDaemonSet(client, maoManagedDaemonSet, framework.MachineAPINamespace, changedDaemonSet)).NotTo(HaveOccurred())

		By(fmt.Sprintf("checking daemonSet %q spec matches", maoManagedDaemonSet))
		Expect(framework.IsDaemonSetSynced(client, initialDaemonSet, maoManagedDaemonSet, framework.MachineAPINamespace)).To(BeTrue())

		By(fmt.Sprintf("checking daemonSet %q is available again", maoManagedDaemonSet))
		Expect(framework.IsDaemonSetAvailable(client, maoManagedDaemonSet, framework.MachineAPINamespace)).To(BeTrue())
	})

	It("reconcile mutating webhook configuration", func() {
		client, err := framework.LoadClient()
		Expect(err).NotTo(HaveOccurred())

		Expect(framework.IsMutatingWebhookConfigurationSynced(client)).To(BeTrue())
	})

	It("reconcile validating webhook configuration", func() {
		client, err := framework.LoadClient()
		Expect(err).NotTo(HaveOccurred())

		Expect(framework.IsValidatingWebhookConfigurationSynced(client)).To(BeTrue())
	})

	It("recover after validating webhook configuration deletion", func() {
		client, err := framework.LoadClient()
		Expect(err).NotTo(HaveOccurred())

		Expect(framework.DeleteValidatingWebhookConfiguration(client, framework.DefaultValidatingWebhookConfiguration)).To(Succeed())

		Expect(framework.IsValidatingWebhookConfigurationSynced(client)).To(BeTrue())
	})

	It("recover after mutating webhook configuration deletion", func() {
		client, err := framework.LoadClient()
		Expect(err).NotTo(HaveOccurred())

		Expect(framework.DeleteMutatingWebhookConfiguration(client, framework.DefaultMutatingWebhookConfiguration)).To(Succeed())

		Expect(framework.IsMutatingWebhookConfigurationSynced(client)).To(BeTrue())
	})

	It("maintains spec after mutating webhook configuration change and preserve caBundle", func() {
		client, err := framework.LoadClient()
		Expect(err).NotTo(HaveOccurred())

		initial, err := framework.GetMutatingWebhookConfiguration(client, framework.DefaultMutatingWebhookConfiguration.Name)
		Expect(initial).ToNot(BeNil())
		Expect(err).To(Succeed())

		toUpdate := initial.DeepCopy()
		for _, webhook := range toUpdate.Webhooks {
			webhook.ClientConfig.CABundle = []byte("test")
			webhook.AdmissionReviewVersions = []string{"test"}
		}

		Expect(framework.UpdateMutatingWebhookConfiguration(client, toUpdate)).To(Succeed())

		Expect(framework.IsMutatingWebhookConfigurationSynced(client)).To(BeTrue())

		updated, err := framework.GetMutatingWebhookConfiguration(client, framework.DefaultMutatingWebhookConfiguration.Name)
		Expect(updated).ToNot(BeNil())
		Expect(err).To(Succeed())

		for i, webhook := range updated.Webhooks {
			Expect(webhook.ClientConfig.CABundle).To(Equal(initial.Webhooks[i].ClientConfig.CABundle))
		}
	})

	It("maintains spec after validating webhook configuration change and preserve caBundle", func() {
		client, err := framework.LoadClient()
		Expect(err).NotTo(HaveOccurred())

		initial, err := framework.GetValidatingWebhookConfiguration(client, framework.DefaultValidatingWebhookConfiguration.Name)
		Expect(initial).ToNot(BeNil())
		Expect(err).To(Succeed())

		toUpdate := initial.DeepCopy()
		for _, webhook := range toUpdate.Webhooks {
			webhook.ClientConfig.CABundle = []byte("test")
			webhook.AdmissionReviewVersions = []string{"test"}
		}

		Expect(framework.UpdateValidatingWebhookConfiguration(client, toUpdate)).To(Succeed())

		Expect(framework.IsValidatingWebhookConfigurationSynced(client)).To(BeTrue())

		updated, err := framework.GetValidatingWebhookConfiguration(client, framework.DefaultValidatingWebhookConfiguration.Name)
		Expect(updated).ToNot(BeNil())
		Expect(err).To(Succeed())

		for i, webhook := range updated.Webhooks {
			Expect(webhook.ClientConfig.CABundle).To(Equal(initial.Webhooks[i].ClientConfig.CABundle))
		}
	})

})

var _ = Describe("[Feature:Operators] Machine API cluster operator status should", func() {
	defer GinkgoRecover()

	It("be available", func() {
		client, err := framework.LoadClient()
		Expect(err).NotTo(HaveOccurred())
		Expect(framework.IsStatusAvailable(client, "machine-api")).To(BeTrue())
	})
})
