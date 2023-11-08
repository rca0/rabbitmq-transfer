package diff

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/dlpco/devops-tools/rabbitmq/loader"
	"github.com/olekukonko/tablewriter"
)

var (
	baseDestUrL   = "https://calm-penguin.rmq.cloudamqp.com"
	baseSourceUrl = "http://172.29.1.9:15672"
)

type jsonBody struct {
	Exchanges   []Exchange
	Queues      []Queue
	Users       []User
	Permissions []Permission
}

type Exchange struct {
	Name   string
	Policy string
	Vhost  string
}

type Queue struct {
	Name  string
	Vhost string
	State string
}

type User struct {
	Name string
	Tag  string
}

type Permission struct {
	Name      string
	Vhost     string
	Configure string
	Write     string
	Read      string
}

func basicAuth(host string) string {
	u, _ := url.Parse(host)
	user := u.User.Username()
	passwd, _ := u.User.Password()
	credentials := fmt.Sprintf("%s:%s", user, passwd)
	return base64.StdEncoding.EncodeToString([]byte(credentials))
}

func getExchanges(amqp, url string) []Exchange {
	var exchanges []Exchange
	var jsonBody jsonBody
	client := &http.Client{}

	url = fmt.Sprintf("%s/api/exchanges", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln("exchange http new request", err)
	}

	req.Header.Add("Authorization", "Basic "+basicAuth(amqp))
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln("exchange client request", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("read all body", err)
	}

	err = json.Unmarshal(body, &jsonBody.Exchanges)
	if err != nil {
		log.Fatalln("unmarshal body", err)
	}

	for _, v := range jsonBody.Exchanges {
		if v.Name == "" || v.Policy == "" || v.Vhost == "" {
			continue
		}
		if strings.Contains(v.Name, "amq") {
			continue
		}
		exchanges = append(exchanges, Exchange{
			Name:   v.Name,
			Policy: v.Policy,
			Vhost:  v.Vhost,
		})
	}
	return exchanges
}

func getQueues(amqp, url string) []Queue {
	var queues []Queue
	var jsonBody jsonBody
	client := &http.Client{}

	url = fmt.Sprintf("%s/api/queues", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln("queues http new request", err)
	}

	req.Header.Add("Authorization", "Basic "+basicAuth(amqp))
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln("queues client request", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("read all body", err)
	}

	err = json.Unmarshal(body, &jsonBody.Queues)
	if err != nil {
		log.Fatalln("unmarshal body", err)
	}
	for _, v := range jsonBody.Queues {
		queues = append(queues, Queue{
			Name:  v.Name,
			Vhost: v.Vhost,
			State: v.State,
		})
	}
	return queues
}

func getUsers(amqp, url string) []User {
	var users []User
	var jsonBody jsonBody
	client := &http.Client{}

	url = fmt.Sprintf("%s/api/users", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln("users http new request", err)
	}

	req.Header.Add("Authorization", "Basic "+basicAuth(amqp))
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln("users client request", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("read all body", err)
	}

	err = json.Unmarshal(body, &jsonBody.Queues)
	if err != nil {
		log.Fatalln("unmarshal body", err)
	}

	for _, v := range jsonBody.Users {
		users = append(users, User{
			Name: v.Name,
		})
	}
	return users
}

func getPermissions(amqp, url string) []Permission {
	var permissions []Permission
	var jsonBody jsonBody
	client := &http.Client{}

	url = fmt.Sprintf("%s/api/permissions", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln("permissions http new request", err)
	}

	req.Header.Add("Authorization", "Basic "+basicAuth(amqp))
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln("permissions client request", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("read all body", err)
	}

	err = json.Unmarshal(body, &jsonBody.Permissions)
	if err != nil {
		log.Fatalln("unmarshal body", err)
	}

	for _, v := range jsonBody.Permissions {
		permissions = append(permissions, Permission{
			Name:      v.Name,
			Vhost:     v.Vhost,
			Configure: v.Configure,
			Write:     v.Write,
			Read:      v.Read,
		})
	}
	return permissions
}

func diffQueues(src, dst []Queue) []Queue {
	var queues []Queue
	srcValue := make(map[string]string)

	for _, s := range dst {
		srcValue[s.Name] = s.Name
	}

	for _, d := range src {
		_, ok := srcValue[d.Name]
		if !ok {
			queues = append(queues, Queue{
				Name:  d.Name,
				Vhost: d.Vhost,
				State: d.State,
			})
		}
	}
	return queues
}

func diffExchanges(src, dst []Exchange) []Exchange {
	var exchanges []Exchange
	srcValue := make(map[string]string)

	for _, s := range dst {
		srcValue[s.Name] = s.Name
	}

	for _, d := range src {
		_, ok := srcValue[d.Name]
		if !ok {
			exchanges = append(exchanges, Exchange{
				Name:   d.Name,
				Policy: d.Policy,
				Vhost:  d.Vhost,
			})
		}
	}
	return exchanges
}

func diffUsers(src, dst []User) []User {
	var users []User
	srcValue := make(map[string]string)

	for _, s := range dst {
		srcValue[s.Name] = s.Name
	}

	for _, d := range src {
		_, ok := srcValue[d.Name]
		if !ok {
			users = append(users, User{
				Name: d.Name,
				Tag:  d.Tag,
			})
		}
	}
	return users
}

func diffPermissions(src, dst []Permission) []Permission {
	var permissions []Permission
	srcValue := make(map[string]string)

	for _, s := range dst {
		srcValue[s.Name] = s.Name
	}

	for _, d := range src {
		_, ok := srcValue[d.Name]
		if !ok {
			permissions = append(permissions, Permission{
				Name:      d.Name,
				Vhost:     d.Vhost,
				Configure: d.Configure,
				Write:     d.Write,
				Read:      d.Read,
			})
		}
	}
	return permissions
}

func printExchangesTable(values []Exchange) {
	var data [][]string

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Exchange", "Policy", "Vhost"})

	for _, v := range values {
		data = append(data, []string{
			v.Name,
			v.Policy,
			v.Vhost,
		})
	}
	table.AppendBulk(data)
	table.Render()
}

func printQueuesTable(values []Queue) {
	var data [][]string

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Vhost", "State"})

	for _, v := range values {
		data = append(data, []string{
			v.Name,
			v.Vhost,
			v.State,
		})
	}
	table.AppendBulk(data)
	table.Render()
}

func printPermissionsTable(values []Permission) {
	var data [][]string

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Vhost", "Configure", "Write", "Read"})

	for _, v := range values {
		data = append(data, []string{
			v.Name,
			v.Vhost,
			v.Configure,
			v.Write,
			v.Read,
		})
	}
	table.AppendBulk(data)
	table.Render()
}

func printUsersTable(values []User) {
	var data [][]string

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Tag"})

	for _, v := range values {
		data = append(data, []string{
			v.Name,
			v.Tag,
		})
	}
	table.AppendBulk(data)
	table.Render()
}

func Run() {
	var c loader.Transfer
	cfg := c.GetTransferConfig()

	srcUsers := getUsers(cfg.Servers.Source, baseSourceUrl)
	dstUsers := getUsers(cfg.Servers.Dest, baseDestUrL)
	log.Printf("Diff: Users from %s to %s servers", baseSourceUrl, baseDestUrL)
	diffu := diffUsers(srcUsers, dstUsers)
	printUsersTable(diffu)

	srcPermissions := getPermissions(cfg.Servers.Source, baseSourceUrl)
	dstPermissions := getPermissions(cfg.Servers.Dest, baseDestUrL)
	log.Printf("Diff: Permissions from %s to %s servers", baseSourceUrl, baseDestUrL)
	diffp := diffPermissions(srcPermissions, dstPermissions)
	printPermissionsTable(diffp)

	srcExchagens := getExchanges(cfg.Servers.Source, baseSourceUrl)
	destExchagens := getExchanges(cfg.Servers.Dest, baseDestUrL)
	log.Printf("Diff: exchanges/vhosts/policies from %s to %s servers", baseSourceUrl, baseDestUrL)
	diffex := diffExchanges(srcExchagens, destExchagens)
	printExchangesTable(diffex)

	srcQueues := getQueues(cfg.Servers.Source, baseSourceUrl)
	dstQueues := getQueues(cfg.Servers.Dest, baseDestUrL)
	log.Printf("Diff: queues/vhosts from %s to %s servers", baseSourceUrl, baseDestUrL)
	diffq := diffQueues(srcQueues, dstQueues)
	printQueuesTable(diffq)
}
