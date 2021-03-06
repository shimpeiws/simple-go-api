package main

import (
	"net/http"
	"io/ioutil"
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type JwtToken struct {
	Jwt string
}

type ItemList struct {
	Data []Item
}

type Item struct {
	Id string
	Attributes ItemDetail
}

type ItemDetail struct {
	Name string
}

type CategoryList struct {
	Data []Item
}

type Category struct {
	Id string
	Attributes CategoryDetail
}

type CategoryDetail struct {
	Name string
}

type CommentList struct {
	Data []Comment
}

type Comment struct {
	Id string
	Attributes CommentDetail
}

type CommentDetail struct {
	Body string
}

type ItemsResponse struct {
	Items ItemList
	Categories CategoryList
	Comments CommentList
}

func getJwtToken() (string, error) {
	requestStr := `{"auth":{"email":"test@example.com","password":"test123"}}`
	req, err := http.NewRequest(
		"POST",
		"http://localhost:3000/user_token",
		bytes.NewBuffer([]byte(requestStr)),
	)

	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	jwtResponse := new(JwtToken)
	fmt.Println(string(body))
	jsonParseError := json.Unmarshal(body, &jwtResponse)
	if jsonParseError != nil {
		return "", jsonParseError
	}

	return jwtResponse.Jwt, nil
}

func getRequest(path, token string) ([]byte, error) {
	req, err := http.NewRequest(
		"GET",
		"http://localhost:3000/" + path,
		nil,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", string("Bearer " + token))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func getItems(token string) (*ItemList, error) {
	body, err := getRequest("items", token)
	if err != nil {
		return nil, err
	}

	itemResponse := new(ItemList)
	fmt.Println("getItemsResponse", string(body))
	jsonParseError := json.Unmarshal(body, &itemResponse)
	if jsonParseError != nil {
		return nil, jsonParseError
	}

	return itemResponse, nil
}

func getCategories(token string) (*CategoryList, error) {
	body, err := getRequest("categories", token)
	if err != nil {
		return nil, err
	}

	categoryResponse := new(CategoryList)
	fmt.Println("categoryResponse", string(body))
	jsonParseError := json.Unmarshal(body, &categoryResponse)
	if jsonParseError != nil {
		return nil, jsonParseError
	}

	return categoryResponse, nil
}

func getComments(token string, itemIds []int) (*CommentList, error) {
	// TODO: implement GET Comments
	commentList := CommentList{}
	return &commentList, nil
}

func main() {
	// Echo instance
	e := echo.New()

	e.Use(middleware.CORS())

	// Route => handler
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!\n")
	})

	e.GET("/users/:id", func(c echo.Context) error {
		// User ID from path `users/:id`
		id := c.Param("id")
		return c.String(http.StatusOK, id)
	})

	e.GET("/rails", func(c echo.Context) error {
		resp, err := http.Get("http://localhost:3000")
		if err != nil {
			return c.String(http.StatusNotFound, "Not Found")
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		return c.String(http.StatusOK, string(body))
	})

	e.GET("/items", func(c echo.Context) error {
		token, err := getJwtToken()
		if err != nil {
			return c.String(http.StatusForbidden, err.Error())
		}
		fmt.Println("token", token)

		// TODO: Do concurrently
		items, getItemErr := getItems(token)
		if getItemErr != nil {
			return c.String(http.StatusNotFound, getItemErr.Error())
		}
		fmt.Println("items", items)

		categories, getCategoryErr := getCategories(token)
		if getCategoryErr != nil {
			return c.String(http.StatusNotFound, getCategoryErr.Error())
		}
		fmt.Println("categories", categories)

		// TODO: implement
		itemIds := []int{1}
		comments, getCommentErr := getComments(token, itemIds)
		if getCommentErr != nil {
			return c.String(http.StatusNotFound, getCommentErr.Error())
		}
		fmt.Println("comments", comments)

		response := ItemsResponse{Items: *items, Categories: *categories, Comments: *comments}
		// response := ItemsResponse{}
		jsonBytes, err := json.Marshal(response)

		fmt.Println("response json", string(jsonBytes))
		return c.String(http.StatusOK, string(jsonBytes))
	})

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
