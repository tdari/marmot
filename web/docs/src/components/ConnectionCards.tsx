import React, { useState } from "react";
import { Icon } from "@iconify/react";
import useDocusaurusContext from "@docusaurus/useDocusaurusContext";

interface Connection {
  name: string;
  description: string;
  href: string;
  icon: string;
  useLocalIcon?: boolean;
}

const connections: Connection[] = [
  {
    name: "Apache Airflow",
    description: "Apache Airflow workflow orchestration connection",
    href: "/docs/Connections/Apache%20Airflow",
    icon: "logos:airflow-icon",
  },
  {
    name: "AWS",
    description: "AWS connection for cloud services authentication",
    href: "/docs/Connections/AWS",
    icon: "logos:amazon-web-services",
  },
  {
    name: "Azure Blob Storage",
    description: "Azure Blob Storage connection",
    href: "/docs/Connections/Azure%20Blob%20Storage",
    icon: "logos:azure-icon",
  },
  {
    name: "ClickHouse",
    description: "ClickHouse columnar database connection",
    href: "/docs/Connections/ClickHouse",
    icon: "simple-icons:clickhouse",
  },
  {
    name: "Google BigQuery",
    description: "Google BigQuery data warehouse connection",
    href: "/docs/Connections/Google%20BigQuery",
    icon: "devicon:googlecloud",
  },
  {
    name: "Google Cloud Storage",
    description: "GCS object storage connection",
    href: "/docs/Connections/Google%20Cloud%20Storage",
    icon: "logos:google-cloud",
  },
  {
    name: "Apache Iceberg",
    description: "Apache Iceberg REST catalog connection",
    href: "/docs/Connections/Apache%20Iceberg",
    icon: "iceberg",
    useLocalIcon: true,
  },
  {
    name: "Apache Kafka",
    description: "Kafka streaming platform connection",
    href: "/docs/Connections/Apache%20Kafka",
    icon: "kafka",
    useLocalIcon: true,
  },
  {
    name: "MongoDB",
    description: "MongoDB NoSQL database connection",
    href: "/docs/Connections/MongoDB",
    icon: "devicon:mongodb",
  },
  {
    name: "MySQL",
    description: "MySQL database connection",
    href: "/docs/Connections/MySQL",
    icon: "devicon:mysql",
  },
  {
    name: "NATS",
    description: "NATS messaging system connection",
    href: "/docs/Connections/NATS",
    icon: "devicon:nats",
  },
  {
    name: "PostgreSQL",
    description: "PostgreSQL database connection",
    href: "/docs/Connections/PostgreSQL",
    icon: "devicon:postgresql",
  },
  {
    name: "Redis",
    description: "Redis in-memory data store connection",
    href: "/docs/Connections/Redis",
    icon: "devicon:redis",
  },
  {
    name: "Trino",
    description: "Trino distributed SQL query engine connection",
    href: "/docs/Connections/Trino",
    icon: "simple-icons:trino",
  },
];

function ConnectionIcon({ connection, isDarkTheme }: { connection: Connection; isDarkTheme: boolean }) {
  if (connection.useLocalIcon) {
    const ext = connection.icon.includes('.') ? '' : '.svg';
    const iconSrc = isDarkTheme ? `/img/dark-${connection.icon}${ext}` : `/img/${connection.icon}${ext}`;
    return <img src={iconSrc} alt={`${connection.name} icon`} className="w-8 h-8" />;
  }

  return (
    <Icon
      icon={connection.icon}
      className="w-8 h-8"
    />
  );
}

export default function ConnectionCards(): JSX.Element {
  const [search, setSearch] = useState("");
  const { siteConfig } = useDocusaurusContext();

  const isDarkTheme = typeof document !== "undefined" &&
    document.documentElement.getAttribute("data-theme") === "dark";

  const filteredConnections = connections.filter(
    (connection) =>
      connection.name.toLowerCase().includes(search.toLowerCase()) ||
      connection.description.toLowerCase().includes(search.toLowerCase())
  );

  return (
    <div className="mt-6">
      <div className="relative mb-5">
        <Icon
          icon="mdi:magnify"
          className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-gray-400 dark:text-gray-500"
        />
        <input
          type="text"
          placeholder="Search connections..."
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          className="w-full pl-9 pr-4 py-2 text-sm bg-white dark:bg-gray-900 border border-gray-200 dark:border-gray-700 rounded-lg text-gray-900 dark:text-white placeholder-gray-400 dark:placeholder-gray-500 focus:outline-none focus:border-[var(--ifm-color-primary)] focus:ring-1 focus:ring-[var(--ifm-color-primary)] transition-colors"
        />
      </div>

      {filteredConnections.length === 0 ? (
        <div className="text-center py-8 text-gray-500 dark:text-gray-400">
          No connections found matching "{search}"
        </div>
      ) : (
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
          {filteredConnections.map((connection) => (
            <a
              key={connection.name}
              href={connection.href}
              className="group block p-5 bg-gray-50 dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700 hover:border-[var(--ifm-color-primary)] dark:hover:border-[var(--ifm-color-primary)] hover:shadow-lg transition-all no-underline"
            >
              <div className="flex items-start gap-4">
                <div className="flex-shrink-0 p-2 bg-white dark:bg-gray-900 rounded-lg border border-gray-200 dark:border-gray-700 group-hover:border-[var(--ifm-color-primary)] transition-colors">
                  <ConnectionIcon connection={connection} isDarkTheme={isDarkTheme} />
                </div>
                <div className="flex-1 min-w-0">
                  <h3 className="text-base font-semibold text-gray-900 dark:text-white m-0 group-hover:text-[var(--ifm-color-primary)] transition-colors">
                    {connection.name}
                  </h3>
                  <p className="mt-1 text-sm text-gray-600 dark:text-gray-400 m-0 line-clamp-2">
                    {connection.description}
                  </p>
                </div>
              </div>
            </a>
          ))}
        </div>
      )}
    </div>
  );
}
