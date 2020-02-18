package controllers

import (
	"context"
	"time"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"

	shipv1beta1 "github.com/tlyng/mockoperator/api/v1beta1"
	component_mocks "github.com/tlyng/mockoperator/component/mocks"
)

var _ = Describe("Frigate Controller Test Suite", func() {
	const timeout = time.Second * 5
	const interval = time.Second * 1

	var (
		mockCtrl        *gomock.Controller
		mockManipulator *component_mocks.MockManipulator
		mgr             ctrl.Manager
		err             error
		stopCh          chan struct{}
		finishCh        chan struct{}
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockManipulator = component_mocks.NewMockManipulator(mockCtrl)
		stopCh = make(chan struct{})
		finishCh = make(chan struct{})

		mgr, err = ctrl.NewManager(cfg, ctrl.Options{
			Scheme: scheme.Scheme,
		})
		Expect(err).ToNot(HaveOccurred())

		err = (&FrigateReconciler{
			Client:      mgr.GetClient(),
			Scheme:      mgr.GetScheme(),
			Log:         ctrl.Log.WithName("controllers").WithName("frigate"),
			Manipulator: mockManipulator,
		}).SetupWithManager(mgr)
		Expect(err).ToNot(HaveOccurred())

		go func() {
			defer close(finishCh)
			defer GinkgoRecover()
			err = mgr.Start(stopCh)
			Expect(err).ToNot(HaveOccurred())
		}()
	})

	AfterEach(func() {
		close(stopCh)
		<-finishCh
		mockCtrl.Finish()
	})

	Describe("Tests", func() {
		var dkey types.NamespacedName

		// BeforeEach(func() {
		// 	By("Creating a frigate")
		// 	dkey = types.NamespacedName{
		// 		Name:      "frigate",
		// 		Namespace: "default",
		// 	}

		// 	frigate := &shipv1beta1.Frigate{
		// 		ObjectMeta: metav1.ObjectMeta{
		// 			Name:      dkey.Name,
		// 			Namespace: dkey.Namespace,
		// 		},
		// 		Spec: shipv1beta1.FrigateSpec{
		// 			Foo: "hello",
		// 		},
		// 	}

		// 	Expect(k8sClient.Create(context.Background(), frigate)).Should(Succeed())

		// 	By("Ensuring frigate is created")
		// 	Eventually(func() error {
		// 		fetched := &shipv1beta1.Frigate{}
		// 		return k8sClient.Get(context.Background(), dkey, fetched)
		// 	}, timeout, interval).Should(Succeed())
		// })

		// AfterEach(func() {
		// 	By("Deleting the frigate")
		// 	delete := &shipv1beta1.Frigate{}
		// 	Expect(k8sClient.Get(context.Background(), dkey, delete)).Should(Succeed())
		// 	Expect(k8sClient.Delete(context.Background(), delete)).Should(Succeed())

		// 	By("Ensuring frigate is deleted")
		// 	Eventually(func() error {
		// 		return k8sClient.Get(context.Background(), dkey, delete)
		// 	}, timeout, interval).ShouldNot(Succeed())
		// })

		It("Frigate Foo should be hello", func() {
			mockManipulator.EXPECT().Manipulate("hello").Return("olleh")
			dkey = types.NamespacedName{
				Name:      "frigate",
				Namespace: "default",
			}

			obj := &shipv1beta1.Frigate{
				ObjectMeta: metav1.ObjectMeta{
					Name:      dkey.Name,
					Namespace: dkey.Namespace,
				},
				Spec: shipv1beta1.FrigateSpec{
					Foo: "hello",
				},
			}

			By("Creating a frigate")
			Expect(k8sClient.Create(context.Background(), obj)).Should(Succeed())
			Eventually(func() string {
				fetched := &shipv1beta1.Frigate{}
				_ = k8sClient.Get(context.Background(), dkey, fetched)
				return fetched.Spec.Foo
			}, timeout, interval).Should(Equal("hello"))

			By("Retrieving the frigate status")
			Eventually(func() string {
				fetched := &shipv1beta1.Frigate{}
				_ = k8sClient.Get(context.Background(), dkey, fetched)
				return fetched.Status.Foo
			}, timeout, interval).Should(Equal("olleh"))
		})
	})
})
