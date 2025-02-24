package discovery

import (
    "fmt"
    "strings"
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
    client          *discovery.DiscoveryClient
    ignoreResources map[string]bool
}

func NewResourceDiscovery(config *rest.Config, ignoreResources []string) (*ResourceDiscovery, error) {
    client, err := discovery.NewDiscoveryClientForConfig(config)
    if err != nil {
        return nil, err
    }

    // Convert ignore list to a map for faster lookups
    ignoreMap := make(map[string]bool)
    for _, r := range ignoreResources {
        ignoreMap[r] = true
    }

    return &ResourceDiscovery{
        client:          client,
        ignoreResources: ignoreMap,
    }, nil
}

func (r *ResourceDiscovery) shouldIgnoreResource(groupVersion, resourceName string) bool {
    // Convert to lowercase for case-insensitive comparison
    resourceName = strings.ToLower(resourceName)
    groupVersion = strings.ToLower(groupVersion)

    // Check exact match
    if r.ignoreResources[resourceName] {
        return true
    }

    // Check with group/version
    fullName := fmt.Sprintf("%s/%s", groupVersion, resourceName)
    if r.ignoreResources[fullName] {
        return true
    }

    return false
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
            // Skip if resource is in ignore list
            if r.shouldIgnoreResource(list.GroupVersion, resource.Name) {
                continue
            }

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