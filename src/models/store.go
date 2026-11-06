package models

import (
	"sync"
	"time"
)

var (
	storeMu sync.RWMutex

	fallbackCategories = []Category{
		{CategoryID: 1, CategoryName: "Portofolio Page"},
		{CategoryID: 2, CategoryName: "ERP"},
		{CategoryID: 3, CategoryName: "Mobile App"},
		{CategoryID: 4, CategoryName: "Games"},
	}

	fallbackClients = []Client{
		{ClientID: 1, FullName: "M Brahmana P", Email: "brahmana@flexoo.id", Username: "brahmana", Password: "password123"},
		{ClientID: 2, FullName: "Demo Client", Email: "demo@flexoo.id", Username: "admin", Password: "admin123"},
	}

	fallbackServices = []Service{
		{
			ServiceID:    1,
			CategoryID:   1,
			CategoryName: "Portofolio Page",
			Name:         "Portofolio Web",
			Description:  "lorem ipsum dolor sit amet consectetur adipiscing elit sed do eiusmod tempor incididunt ut labore et.",
			PageRange:    "1-3 Pages",
			BasePrice:    199.00,
		},
		{
			ServiceID:    2,
			CategoryID:   3,
			CategoryName: "Mobile App",
			Name:         "Mobile App",
			Description:  "Custom cross-platform mobile application built for iOS and Android with modern UI design and scalable cloud backend.",
			PageRange:    "1-4 Pages",
			BasePrice:    299.00,
		},
		{
			ServiceID:    3,
			CategoryID:   2,
			CategoryName: "ERP",
			Name:         "ERP",
			Description:  "Comprehensive Enterprise Resource Planning platform tailored to automate business operations, finance, inventory, and reporting.",
			PageRange:    "5-20 Pages",
			BasePrice:    9999.00,
		},
		{
			ServiceID:    4,
			CategoryID:   4,
			CategoryName: "Games",
			Name:         "Custom",
			Description:  "Chat with customer support for custom package inquiries.",
			PageRange:    "Chat with customer support for custom package inquiries.",
			BasePrice:    9999.00,
		},
	}

	fallbackDeliverables = []ServiceDeliverable{
		{DeliverableID: 1, ServiceID: 1, Label: "HTML, CSS, JS Files"},
		{DeliverableID: 2, ServiceID: 1, Label: "Free Vercel Domain"},
		{DeliverableID: 3, ServiceID: 1, Label: "Static Hosting"},
		{DeliverableID: 4, ServiceID: 2, Label: "iOS & Android Builds"},
		{DeliverableID: 5, ServiceID: 2, Label: "API & Push Notifications"},
		{DeliverableID: 6, ServiceID: 2, Label: "Store Deployment Guide"},
		{DeliverableID: 7, ServiceID: 3, Label: "Full-Stack Source Code"},
		{DeliverableID: 8, ServiceID: 3, Label: "Neon PostgreSQL Database"},
		{DeliverableID: 9, ServiceID: 3, Label: "Role-Based Access Control"},
		{DeliverableID: 10, ServiceID: 4, Label: "Custom System Architecture"},
		{DeliverableID: 11, ServiceID: 4, Label: "Dedicated Support Channel"},
		{DeliverableID: 12, ServiceID: 4, Label: "Flexible Milestone Delivery"},
	}

	fallbackPastWorks = []PastWork{
		{WorkID: 1, ServiceID: 1, Title: "SaaS Platform", ImageURL: "/static/img/saas-preview.png"},
		{WorkID: 2, ServiceID: 2, Title: "SaaS Platform", ImageURL: "/static/img/saas-preview.png"},
		{WorkID: 3, ServiceID: 3, Title: "SaaS Platform", ImageURL: "/static/img/saas-preview.png"},
		{WorkID: 4, ServiceID: 4, Title: "SaaS Platform", ImageURL: "/static/img/saas-preview.png"},
	}

	fallbackOrders = []ServiceOrder{
		{
			OrderID:           "ORDER-5123",
			ClientID:          1,
			ServiceID:         1,
			ServiceName:       "Portofolio Web",
			ProjectName:       "HIMASAKTA ITS WEB",
			Scale:             "Standard",
			Comments:          "Official web portal for HIMASAKTA ITS",
			SRSFileURL:        "SRS_HIMASAKTA_ITS.pdf",
			Status:            "In Progress",
			EstimatedDelivery: time.Date(2026, 11, 18, 0, 0, 0, 0, time.UTC),
			CreatedAt:         time.Date(2026, 11, 7, 2, 0, 0, 0, time.UTC),
		},
	}
)

