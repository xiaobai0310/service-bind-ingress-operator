package utils

import (
	"bytes"
	"github.com/xiaobai0310/service-bind-ingress-operator/api/v1alpha1"
	"gopkg.in/yaml.v2"
	appv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkv1 "k8s.io/api/networking/v1"
	"text/template"
)

func parseTemplate(templateName string, instance *v1alpha1.ServiceBindingIngress) []byte {
	tmpl, err := template.ParseFiles("controller/templates/" + templateName + ".yml")
	if err != nil {
		panic(err)
	}
	b := new(bytes.Buffer)
	err = tmpl.Execute(b, instance)
	if err != nil {
		panic(err)
	}
	return b.Bytes()
}

func NewDeployment(instance *v1alpha1.ServiceBindingIngress) *appv1.Deployment {
	d := &appv1.Deployment{}
	err := yaml.Unmarshal(parseTemplate("deployment", instance), d)
	if err != nil {
		panic(err)
	}
	return d
}

func NewIngress(instance *v1alpha1.ServiceBindingIngress) *networkv1.Ingress {
	i := &networkv1.Ingress{}
	err := yaml.Unmarshal(parseTemplate("ingress", instance), i)
	if err != nil {
		panic(err)
	}
	return i
}

func NewService(instance *v1alpha1.ServiceBindingIngress) *corev1.Service {
	s := &corev1.Service{}
	err := yaml.Unmarshal(parseTemplate("service", instance), s)
	if err != nil {
		panic(err)
	}
	return s
}
