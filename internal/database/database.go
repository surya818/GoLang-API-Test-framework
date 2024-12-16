// Copyright © 2024 Kong Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package database

import (
	"database/sql"
	"fmt"
	"log"
)

// NewDatabase creates a new SQLite database instance with pre-populated data.
func NewDatabase() (*sql.DB, error) {
	// Connect to the SQLite database
	db, err := sql.Open("sqlite3", "./candidate-take-home-exercise-sdet.db")
	if err != nil {
		return nil, fmt.Errorf("unable to create andidate-take-home-exercise-sdet.db: %w", err)
	}

	// Drop existing tables to ensure a clean start
	_, err = db.Exec(`DROP TABLE IF EXISTS service_versions`)
	if err != nil {
		return nil, fmt.Errorf("unable to DROP service_versions table: %w", err)
	}
	_, err = db.Exec(`DROP TABLE IF EXISTS services`)
	if err != nil {
		return nil, fmt.Errorf("unable to DROP services table: %w", err)
	}

	// Create the tables
	_, err = db.Exec(`
        CREATE TABLE services (
            id TEXT PRIMARY KEY,
            name TEXT CHECK(length(name) <= 64),
            description TEXT CHECK(length(description) <= 255),
            created_at DATETIME NOT NULL,
            updated_at DATETIME NOT NULL
        )
    `)
	if err != nil {
		return nil, fmt.Errorf("unable to CREATE services table: %w", err)
	}
	_, err = db.Exec(`
        CREATE TABLE service_versions (
            id TEXT PRIMARY KEY,
            service_id TEXT NOT NULL,
            version TEXT NOT NULL CHECK(length(version) <= 16),
            created_at DATETIME NOT NULL,
            updated_at DATETIME NOT NULL,
            FOREIGN KEY (service_id) REFERENCES services(id)
        )
    `)
	if err != nil {
		return nil, fmt.Errorf("unable to CREATE service_versions table: %w", err)
	}

	// Insert initial data to populate the services table
	//nolint:lll
	initialServices := []struct {
		ID, Name, Description string
	}{
		{"01836a4b-c000-7fd0-b89a-c0e51546b001", "Locate Us", "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Turpis non a, pellentesque ipsumaliquet id ligula, tincidunt ut orci."},
		{"01836a4b-c001-7fd1-b89a-c0e51546b002", "Collect Money", ""},
		{"01836a4b-c002-7fd2-b89a-c0e51546b003", "Contact Us", "Lorem ipsum dolor sit amet, consectetur adipiscing"},
		{"01836a4b-c003-7fd3-b89a-c0e51546b004", "Contact Us", "Lorem ipsum dolor sit amet, consectetur adipiscing"},
		{"01836a4b-c004-7fd4-b89a-c0e51546b005", "FX Rates International", "Lorem ipsum dolor"},
		{"01836a4bc0057fd5b89ac0e51546b006", "FX Rates International [Internal]", "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Turpis non a, pellentesque ipsumaliquet id, ornare in lacus."},
		{"01836a4b-c006-7fd6-b89a-c0e51546b007", "Notifications", ""},
		{"01836a4b-c007-7fd7-b89a-c0e51546b008", "Notifications", ""},
		{"01836a4b-c008-7fd8-b89a-c0e51546b009", "Priority Services", ""},
		{"01836a4b-c009-7fd9-b89a-c0e51546b00a", "Reporting", "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Turpis non a, pellentesque ipsumaliquet id arcu, gravida quis."},
		{"01836a4b-c00a-7fda-b89a-c0e51546b00b", "Security", "Lorem ipsum dolor"},
		{"01836a4b-c00b-7fdb-b89a-c0e51546b00c", "Security", "Lorem ipsum dolor"},
	}

	for _, service := range initialServices {
		_, err := db.Exec(`
			INSERT INTO services (id, name, description, created_at, updated_at)
			VALUES (?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		`, service.ID, service.Name, service.Description)
		if err != nil {
			return nil, fmt.Errorf("unable to INSERT service %s into services table: %w", service.Name, err)
		}
	}

	// Insert initial data to populate the service_versions table
	initialVersions := []struct {
		ID, ServiceID, Version string
	}{
		// Locate Us Service Versions
		{"01836b4b-c100-7fd0-b89a-c0e51546b101", "01836a4b-c000-7fd0-b89a-c0e51546b001", "v1.0"},
		{"01836b4b-c101-7fd1-b89a-c0e51546b102", "01836a4b-c000-7fd0-b89a-c0e51546b001", "v1.1-alpha"},
		{"01836b4b-c102-7fd2-b89a-c0e51546b103", "01836a4b-c000-7fd0-b89a-c0e51546b001", "v2.0-beta"},

		// Collect Money Service Versions
		{"01836b4b-c103-7fd3-b89a-c0e51546b104", "01836a4b-c001-7fd1-b89a-c0e51546b002", "v1.0"},
		{"01836b4b-c104-7fd4-b89a-c0e51546b105", "01836a4b-c001-7fd1-b89a-c0e51546b002", "v1.2-preview"},
		{"01836b4b-c105-7fd5-b89a-c0e51546b106", "01836a4b-c001-7fd1-b89a-c0e51546b002", "v2.1-stable"},

		// Contact Us Service Versions
		{"01836b4b-c106-7fd6-b89a-c0e51546b107", "01836a4b-c002-7fd2-b89a-c0e51546b003", "v1.0"},
		{"01836b4b-c107-7fd7-b89a-c0e51546b108", "01836a4b-c002-7fd2-b89a-c0e51546b003", "v1.5"},
		{"01836b4b-c108-7fd8-b89a-c0e51546b109", "01836a4b-c002-7fd2-b89a-c0e51546b003", "v2.3"},

		// Contact Us (Duplicate) Service Versions
		{"01836b4b-c109-7fd9-b89a-c0e51546b110", "01836a4b-c003-7fd3-b89a-c0e51546b004", "v1.0"},
		{"01836b4b-c110-7fda-b89a-c0e51546b111", "01836a4b-c003-7fd3-b89a-c0e51546b004", "v1.2"},
		{"01836b4b-c111-7fdb-b89a-c0e51546b112", "01836a4b-c003-7fd3-b89a-c0e51546b004", "v2.0"},

		// FX Rates International Service Versions
		{"01836b4b-c112-7fdc-b89a-c0e51546b113", "01836a4b-c004-7fd4-b89a-c0e51546b005", "v1.0"},
		{"01836b4b-c113-7fdd-b89a-c0e51546b114", "01836a4b-c004-7fd4-b89a-c0e51546b005", "v1.3-staging"},
		{"01836b4b-c114-7fde-b89a-c0e51546b115", "01836a4b-c004-7fd4-b89a-c0e51546b005", "v2.0-release"},

		// FX Rates International (Extended) Service Versions
		{"01836b4b-c115-7fdf-b89a-c0e51546b116", "01836a4bc0057fd5b89ac0e51546b006", "v1.0"},
		{"01836b4b-c116-7fe0-b89a-c0e51546b117", "01836a4bc0057fd5b89ac0e51546b006", "v1.2"},
		{"01836b4b-c117-7fe1-b89a-c0e51546b118", "01836a4bc0057fd5b89ac0e51546b006", "v2.1"},

		// Notifications Service Versions
		{"01836b4b-c106-7fd6-b89a-c0e51546b119", "01836a4b-c006-7fd6-b89a-c0e51546b007", "v1.1"},
		{"01836b4b-c107-7fd7-b89a-c0e51546b120", "01836a4b-c006-7fd6-b89a-c0e51546b007", "v1.8"},
		{"01836b4b-c108-7fd8-b89a-c0e51546b121", "01836a4b-c006-7fd6-b89a-c0e51546b007", "v7.3"},

		// Notifications (Duplicate) Service Versions
		{"01836b4b-c109-7fd9-b89a-c0e51546b122", "01836a4b-c007-7fd7-b89a-c0e51546b008", "v5.0"},
		{"01836b4b-c110-7fda-b89a-c0e51546b123", "01836a4b-c007-7fd7-b89a-c0e51546b008", "v4.2"},
		{"01836b4b-c111-7fdb-b89a-c0e51546b124", "01836a4b-c007-7fd7-b89a-c0e51546b008", "v9.0"},

		// Priority Services Versions
		{"01836b4b-c118-7fe2-b89a-c0e51546b125", "01836a4b-c008-7fd8-b89a-c0e51546b009", "v1.0"},
		{"01836b4b-c119-7fe3-b89a-c0e51546b126", "01836a4b-c008-7fd8-b89a-c0e51546b009", "v1.4"},
		{"01836b4b-c120-7fe4-b89a-c0e51546b127", "01836a4b-c008-7fd8-b89a-c0e51546b009", "v2.0"},

		// Reporting Versions
		{"01836b4b-c121-7fe5-b89a-c0e51546b128", "01836a4b-c009-7fd9-b89a-c0e51546b00a", "v1.0"},
		{"01836b4b-c122-7fe6-b89a-c0e51546b129", "01836a4b-c009-7fd9-b89a-c0e51546b00a", "v1.3"},
		{"01836b4b-c123-7fe7-b89a-c0e51546b130", "01836a4b-c009-7fd9-b89a-c0e51546b00a", "v2.0"},

		// Security Versions
		{"01836b4b-c124-7fe8-b89a-c0e51546b131", "01836a4b-c00a-7fda-b89a-c0e51546b00b", "v1.0"},
		{"01836b4b-c125-7fe9-b89a-c0e51546b132", "01836a4b-c00a-7fda-b89a-c0e51546b00b", "v1.2"},
		{"01836b4b-c126-7fea-b89a-c0e51546b133", "01836a4b-c00a-7fda-b89a-c0e51546b00b", "v2.0"},

		// Security (Duplicate) Versions
		{"01836b4b-c127-7feb-b89a-c0e51546b134", "01836a4b-c00b-7fdb-b89a-c0e51546b00c", "v1.0"},
		{"01836b4b-c128-7fec-b89a-c0e51546b135", "01836a4b-c00b-7fdb-b89a-c0e51546b00c", "v1.1"},
		{"01836b4b-c129-7fed-b89a-c0e51546b136", "01836a4b-c00b-7fdb-b89a-c0e51546b00c", "v2.0"},
	}

	for _, version := range initialVersions {
		_, err := db.Exec(`
            INSERT INTO service_versions (id, service_id, version, created_at, updated_at)
            VALUES (?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        `, version.ID, version.ServiceID, version.Version)
		if err != nil {
			log.Printf("Failed to insert version %s for service %s: %v", version.Version, version.ServiceID, err)
			return nil, fmt.Errorf("unable to INSERT version %s for service %s into versions table: %w", version.Version,
				version.ServiceID, err)
		}
	}

	return db, nil
}
