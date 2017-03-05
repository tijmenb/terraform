package kubernetes

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"k8s.io/kubernetes/pkg/api/errors"
	api "k8s.io/kubernetes/pkg/api/v1"
	kubernetes "k8s.io/kubernetes/pkg/client/clientset_generated/release_1_5"
)

func resourceKubernetesConfigMap() *schema.Resource {
	return &schema.Resource{
		Create: resourceKubernetesConfigMapCreate,
		Read:   resourceKubernetesConfigMapRead,
		Update: resourceKubernetesConfigMapUpdate,
		Delete: resourceKubernetesConfigMapDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"metadata": metadataSchema,
			"data": {
				Type:     schema.TypeMap,
				Optional: true,
			},
		},
	}
}

func resourceKubernetesConfigMapCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*kubernetes.Clientset)

	metadata := expandMetadata(d.Get("metadata").([]interface{}))
	cfgMap := api.ConfigMap{
		ObjectMeta: metadata,
		Data:       d.Get("data").(map[string]string),
	}
	log.Printf("[INFO] Creating new config map: %#v", cfgMap)
	out, err := conn.CoreV1().ConfigMaps(metadata.Namespace).Create(&cfgMap)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Submitted new config map: %#v", out)
	d.SetId(out.Name)

	return resourceKubernetesConfigMapRead(d, meta)
}

func resourceKubernetesConfigMapRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*kubernetes.Clientset)

	metadata := expandMetadata(d.Get("metadata").([]interface{}))
	name := d.Id()
	log.Printf("[INFO] Reading config map %s", name)
	cfgMap, err := conn.CoreV1().ConfigMaps(metadata.Namespace).Get(name)
	if err != nil {
		if statusErr, ok := err.(*errors.StatusError); ok && statusErr.ErrStatus.Code == 404 {
			log.Printf("[WARN] Removing config map %s (it is gone)", name)
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] Received error: %#v", err)
		return err
	}
	log.Printf("[INFO] Received config map: %#v", cfgMap)
	err = d.Set("metadata", flattenMetadata(cfgMap.ObjectMeta))
	if err != nil {
		return err
	}
	d.Set("data", cfgMap.Data)

	return nil
}

func resourceKubernetesConfigMapUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*kubernetes.Clientset)

	metadata := expandMetadata(d.Get("metadata").([]interface{}))
	// This is necessary in case the name is generated
	metadata.Name = d.Id()

	cfgMap := api.ConfigMap{
		ObjectMeta: metadata,
	}
	log.Printf("[INFO] Updating config map: %#v", cfgMap)
	out, err := conn.CoreV1().ConfigMaps(metadata.Namespace).Update(&cfgMap)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Submitted updated config map: %#v", out)
	d.SetId(out.Name)

	return resourceKubernetesConfigMapRead(d, meta)
}

func resourceKubernetesConfigMapDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*kubernetes.Clientset)

	metadata := expandMetadata(d.Get("metadata").([]interface{}))
	name := d.Id()
	log.Printf("[INFO] Deleting config map: %#v", name)
	err := conn.CoreV1().ConfigMaps(metadata.Namespace).Delete(name, &api.DeleteOptions{})
	if err != nil {
		return err
	}

	log.Printf("[INFO] Config map %s deleted", name)

	d.SetId("")
	return nil
}
