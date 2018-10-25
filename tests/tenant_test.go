package tests

import (
	"fmt"
	"testing"

	"github.com/ciscoecosystem/aci-go-client/client"
	"github.com/ciscoecosystem/aci-go-client/models"
)

// func GetTestClient() *client.Client {
// 	return client.GetClient("https://192.168.10.102", "admin", client.Insecure(true), client.PrivateKey("C:\\Users\\Crest\\Desktop\\certtest\\admin.key"), client.AdminCert("test.crt"))

// }

func GetTestClient() *client.Client {
	return client.GetClient("https://192.168.10.102", "admin", client.Insecure(true), client.Password("cisco123"))

}

// func TestTenantPrepareModel(t *testing.T) {
// 	c := GetTestClient()

// 	cont, _, err := c.PrepareModel(models.NewTenant("terraform-test-tenant", "A test tenant created with aci-client-sdk."))

// 	if err != nil {
// 		t.Error(err)
// 	}
// 	if !cont.ExistsP("FvTenant.attributes.dn") {
// 		t.Error("malformed model")
// 	}
// }

func createTenant(c *client.Client, dn string, desc string) (*models.Tenant, error) {
	tenant := models.NewTenant(fmt.Sprintf("tn-%s", dn), "uni", desc, models.TenantAttributes{})
	err := c.Save(tenant)
	return tenant, err
}

func deleteTenant(c *client.Client, tenant *models.Tenant) error {
	return c.Delete(tenant)
}

func TestTenantCreation(t *testing.T) {
	c := GetTestClient()
	tenant, err := createTenant(c, "terraform-test-tenant", "A test tenant created with aci-client-sdk.")

	if err != nil {
		t.Error(err)
	}

	err = deleteTenant(c, tenant)
	if err != nil {
		t.Error(err)
	}
}

func TestDuplicateTenant(t *testing.T) {
	c := GetTestClient()
	tenant1, err := createTenant(c, "terraform-test-tenant", "A test tenant created with aci-client-sdk.")
	if err != nil {
		t.Error(err)
	}
	_, err = createTenant(c, "terraform-test-tenant", "A test tenant created with aci-client-sdk.")
	if err != nil {
		t.Error(err)
	}

	err = deleteTenant(c, tenant1)
	if err != nil {
		t.Error(err)
	}

}

func TestGetTenantContainer(t *testing.T) {

	c := GetTestClient()
	tenant, _ := createTenant(c, "terraform-test-tenant", "A test tenant created with aci-client-sdk.")
	cont, err := c.Get("uni/tn-terraform-test-tenant")

	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", cont)

	err = deleteTenant(c, tenant)
	if err != nil {
		t.Error(err)
	}
}

func TestTenantFromContainer(t *testing.T) {
	c := GetTestClient()
	tenant, _ := createTenant(c, "terraform-test-tenant", "A test tenant created with aci-client-sdk.")
	cont, err := c.Get("uni/tn-terraform-test-tenant")
	if err != nil {
		t.Error(err)
	}
	tenantCon := models.TenantFromContainer(cont)
	fmt.Println(tenantCon.DistinguishedName)
	if tenantCon.DistinguishedName == "" {
		t.Error("the tenant dn was empty")
	}
	err = deleteTenant(c, tenant)
	if err != nil {
		t.Error(err)
	}
}

func TestUpdateTenant(t *testing.T) {
	client := GetTestClient()
	tenant, _ := createTenant(client, "terraform-test-tenant", "A test tenant created with aci-client-sdk.")
	cont, err := client.Get("uni/tn-terraform-test-tenant")
	if err != nil {
		t.Error(err)
	}
	tenantCon := models.TenantFromContainer(cont)
	if tenantCon.DistinguishedName == "" {
		t.Error("the tenant dn was empty")
	}
	tenantCon.Description = "Updated the description "
	err = client.Save(tenantCon)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("Description Updated for tenant")
	err = deleteTenant(client, tenant)
	if err != nil {
		t.Error(err)
	}
}

func TestTenantDelete(t *testing.T) {
	c := GetTestClient()
	tenant, _ := createTenant(c, "terraform-test-tenant", "A test tenant created with aci-client-sdk.")
	cont, err := c.Get("uni/tn-terraform-test-tenant")
	if err != nil {
		t.Error(err)
	}
	tenantCon := models.TenantFromContainer(cont)
	fmt.Println(tenantCon.DistinguishedName)
	if tenantCon.DistinguishedName == "" {
		t.Error("the tenant dn was empty")
	}

	err = c.Delete(tenant)
	if err != nil {
		t.Error("the tenant was not remove")
	}
	err = deleteTenant(c, tenant)
	if err != nil {
		t.Error(err)
	}

}

func TestCreateRelationToVzFilter(t *testing.T) {
	c := GetTestClient()
	_, _ = createTenant(c, "terraform-test-tenant", "A test tenant created with aci-client-sdk.")
	cont, err := c.Get("uni/tn-terraform-test-tenant")
	if err != nil {
		t.Error(err)
	}
	tenantCon := models.TenantFromContainer(cont)
	fmt.Println(tenantCon.DistinguishedName)
	err = c.CreateRelationTovzFilter("uni/tn-terraform-test-tenant", "uni/tn-terraform-test-tenant/flt-test-rel112")
	if err != nil {
		t.Error(err)
	}

}

func TestDeleteRelationToVzFilter(t *testing.T) {
	c := GetTestClient()
	_, _ = createTenant(c, "terraform-test-tenant", "A test tenant created with aci-client-sdk.")
	cont, err := c.Get("uni/tn-terraform-test-tenant")
	if err != nil {
		t.Error(err)
	}
	tenantCon := models.TenantFromContainer(cont)
	fmt.Println(tenantCon.DistinguishedName)
	err = c.DeleteRelationTovzFilter("uni/tn-terraform-test-tenant", "uni/tn-terraform-test-tenant/flt-test-rel112")
	if err != nil {
		t.Error(err)
	}

}
