package main

import (
	"context"
	"fmt"

	fisalisgo "github.com/LillySchramm/FiSaLisGo.git"
)

func main() {
	ctx := context.Background()

	res, err := fisalisgo.Search(ctx, "Vladimir Vladimirovich Putin")
	if err != nil {
		panic(err)
	}

	for _, r := range res {
		println(fmt.Sprintf("ID: %s", r.Id))
		println(fmt.Sprintf("Match: %f", r.Match))
		println(fmt.Sprintf("Description: %s", r.Description))

		for _, d := range r.Documents {
			println(fmt.Sprintf("Document ID: %s", d.Id))
			println(fmt.Sprintf("Document Link: %s", d.Link))
			println(fmt.Sprintf("Document Date: %s", d.Date.String()))

			println("Names:")
			for _, n := range d.Names {
				println(fmt.Sprintf("  %s", n))
			}

			println("BirthYears:")
			for _, b := range d.BirthDates {
				println(fmt.Sprintf("  %s", b.String()))
			}

			println("Orgs:")
			for _, o := range d.Orgs {
				println(fmt.Sprintf("  %s", o))
			}
		}

		println("-------------------")
	}

}
