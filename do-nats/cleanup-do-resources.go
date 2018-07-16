package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"
)

var doPat *string
var tag *string
var client *godo.Client

func main() {

	doPat = flag.String("DO_PERSONAL_ACCESS_TOKEN", "", "DigitalOcean Personal Access Token")
	tag = flag.String("DROPLET_TAG", "", "Droplet tag")

	flag.Parse()

	if *doPat == "" {
		panic("Flag DO_PERSONAL_ACCESS_TOKEN missing!")
	}

	if *tag == "" {
		panic("Flag DROPLET_TAG missing!")
	}
	client = getDOClient()
	cleanUp()
}

func cleanUp() {
	lbID := getLoadBalancerID()
	removeLoadBalancer(lbID)

	removeDropletsWithTag()
}

func removeLoadBalancer(lbID string) {
	fmt.Println("Deleting LoadBalancer with ID " + lbID)
	resp, err := client.LoadBalancers.Delete(context.TODO(), lbID)
	fmt.Println("LoadBalancer removal API response - " + resp.Status)
	if err != nil {
		panic("Could not remove LoadBalancer " + lbID + " due to " + err.Error())
	} else {
		fmt.Println("Deleted LoadBalancer with ID " + lbID)
	}
}

func removeDropletsWithTag() {
	fmt.Println("Deleting Droplets wit Tag - " + *tag)
	resp, err := client.Droplets.DeleteByTag(context.TODO(), *tag)
	fmt.Println("Droplet removal API response - " + resp.Status)
	if err != nil {
		panic("Could not remove droplets with tag " + *tag + " due to " + err.Error())
	} else {
		fmt.Println("Deleted Droplets with tag " + *tag)
	}
}

//assuming that we have only ONE load balancer
func getLoadBalancerID() string {
	opt := &godo.ListOptions{
		Page:    1,
		PerPage: 200,
	}
	ctx := context.TODO()
	lbs, _, err := client.LoadBalancers.List(ctx, opt)
	if err != nil {
		panic("Could not find Load Balancers due to " + err.Error())
	}

	if len(lbs) < 1 {
		panic("No LoadBalancers found")
	}
	lb := lbs[0]
	fmt.Println("Got LoadBalancer ID - " + lb.ID)
	return lb.ID
}

type tokenSource struct {
	accessToken string
}

func (t *tokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.accessToken,
	}
	return token, nil
}

func getDOClient() *godo.Client {
	tokenSource := &tokenSource{
		accessToken: *doPat,
	}

	oauthClient := oauth2.NewClient(oauth2.NoContext, tokenSource)
	client := godo.NewClient(oauthClient)
	return client
}
