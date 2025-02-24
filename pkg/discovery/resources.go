package discovery

import (
    "k8s.io/client-go/discovery"
    "k8s.io/client-go/rest"
    "k8s.io/apimachinery/pkg/runtime/schema"
)

type APIResource struct {
    Group      string
    Name       string
    Namespaced bool
    GVR        schema.GroupVersionResource
}

type ResourceDiscovery struct {
    client *discovery.DiscoveryClient
}

func NewResourceDiscovery(config *rest.Config) (*ResourceDiscovery, error) {
    client, err := discovery.NewDiscoveryClientForConfig(config)
    if err != nil {
        return nil, err
    }
    return &ResourceDiscovery{client: client}, nil
}

func (r *ResourceDiscovery) GetAPIResources() ([]APIResource, error) {
    lists, err := r.client.ServerPreferredResources()
    if err != nil {
        return nil, err
    }

    var resources []APIResource
    for _, list := range lists {
        gv, err := schema.ParseGroupVersion(list.GroupVersion)
        if err != nil {
            continue
        }

        for _, resource := range list.APIResources {
            // Only include resources that support GET/LIST operations
            if containsString(resource.Verbs, "list") && containsString(resource.Verbs, "get") {
                gvr := schema.GroupVersionResource{
                    Group:    gv.Group,
                    Version:  gv.Version,
                    Resource: resource.Name,
                }
                
                resources = append(resources, APIResource{
                    Group:      list.GroupVersion,
                    Name:       resource.Name,
                    Namespaced: resource.Namespaced,
                    GVR:        gvr,
                })
            }
        }
    }
    return resources, nil
}

func containsString(slice []string, s string) bool {
    for _, item := range slice {
        if item == s {
            return true
        }
    }
    return false
} 