import { useEffect, useState } from "react";
import axios from "axios";
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  Tooltip,
  ResponsiveContainer,
} from "recharts";

export default function App() {
  const [deps, setDeps] = useState([]);
  const [name, setName] = useState("");
  const [search, setSearch] = useState("");
  const [loading, setLoading] = useState(false);

  const loadAll = async () => {
    try {
      setLoading(true);
      const res = await axios.get(
        "http://localhost:8080/api/v1/dependencies"
      );
      setDeps(res.data);
    } catch (err) {
      console.error("Failed to load dependencies", err);
    } finally {
      setLoading(false);
    }
  };

  const searchByName = async () => {
    if (!search.trim()) {
      loadAll();
      return;
    }

    try {
      setLoading(true);
      const res = await axios.get(
        `http://localhost:8080/api/v1/dependencies/${search}`
      );

      setDeps(res.data);
    } catch (err) {
      console.error("Search failed", err);
      setDeps([]);
    } finally {
      setLoading(false);
    }
  };

  const sync = async () => {
    if (!name.trim()) return;

    try {
      setLoading(true);
      await axios.put(
        `http://localhost:8080/api/v1/dependencies/${name}`
      );
      setName("");
      await loadAll();
    } catch (err) {
      console.error("Sync failed", err);
    } finally {
      setLoading(false);
    }
  };

  const remove = async () => {
    if (!name.trim()) return;

    try {
      setLoading(true);
      await axios.delete(
        `http://localhost:8080/api/v1/dependencies/${name}`
      );
      setName("");
      await loadAll();
    } catch (err) {
      console.error("Delete failed", err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadAll();
  }, []);

  return (
    <div style={styles.outer}>
      <div style={styles.page}>
        <h1 style={styles.title}>Dependency Dashboard</h1>

        {/* -------- Manage Section -------- */}
        <section style={styles.card}>
          <h2>Manage package</h2>

          <div style={styles.column}>
            <input
              style={styles.input}
              placeholder="e.g. react"
              value={name}
              onChange={(e) => setName(e.target.value)}
            />

            <button style={styles.buttonPrimary} onClick={sync}>
              Sync Package
            </button>

            <button style={styles.buttonDanger} onClick={remove}>
              Delete Package
            </button>
          </div>
        </section>

        {/* -------- Search Section -------- */}
        <section style={styles.card}>
          <h2>Search by name</h2>

          <div style={styles.column}>
            <input
              style={styles.input}
              placeholder="Exact package by name..."
              value={search}
              onChange={(e) => setSearch(e.target.value)}
            />

            <button style={styles.buttonPrimary} onClick={searchByName}>
              Search
            </button>

            <button style={styles.buttonSecondary} onClick={loadAll}>
              Clear Search
            </button>
          </div>
        </section>

        {/* -------- Table Section -------- */}
        <section style={styles.card}>
          <h2>Stored Dependencies</h2>

          {loading && <p>Loading...</p>}

          {!loading && deps.length === 0 && (
            <p>No dependencies found.</p>
          )}

          <div style={{ overflowX: "auto" }}>
            <table style={styles.table}>
              <thead>
                <tr>
                  <th>Name</th>
                  <th>Version</th>
                  <th>Score</th>
                  <th>Last Updated</th>
                </tr>
              </thead>
              <tbody>
                {deps.map((d) => (
                  <tr key={d.name}>
                    <td>{d.name}</td>
                    <td>{d.version}</td>
                    <td>{d.openssf_score}</td>
                    <td>
                      {new Date(d.last_updated).toLocaleString()}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </section>

        {/* -------- Chart Section -------- */}
        <section style={styles.card}>
          <h2>OpenSSF Score Overview</h2>

          <ResponsiveContainer width="100%" height={300}>
            <BarChart data={deps}>
              <XAxis dataKey="name" />
              <YAxis />
              <Tooltip />
              <Bar dataKey="openssf_score" />
            </BarChart>
          </ResponsiveContainer>
        </section>
      </div>
    </div>
  );
}

/* ---------- Styling ---------- */

const styles = {
  outer: {
    minHeight: "100vh",
    width: "100%",
    display: "flex",
    justifyContent: "center",
    alignItems: "center",
    backgroundColor: "#f4f6f8",
    padding: "20px",
  },
  page: {
    width: "100%",
    maxWidth: "1000px",
    fontFamily: "Arial, sans-serif",
  },
  title: {
    marginBottom: "30px",
  },
  card: {
    backgroundColor: "white",
    padding: "20px",
    marginBottom: "30px",
    borderRadius: "8px",
    boxShadow: "0 2px 6px rgba(0,0,0,0.08)",
  },
  column: {
    display: "flex",
    flexDirection: "column",
    gap: "10px",
  },
  input: {
    padding: "8px",
    fontSize: "14px",
  },
  buttonPrimary: {
    padding: "8px 16px",
    cursor: "pointer",
    backgroundColor: "#2563eb",
    color: "white",
    border: "none",
    borderRadius: "4px",
  },
  buttonDanger: {
    padding: "8px 16px",
    cursor: "pointer",
    backgroundColor: "#dc2626",
    color: "white",
    border: "none",
    borderRadius: "4px",
  },
  buttonSecondary: {
    padding: "8px 16px",
    cursor: "pointer",
    backgroundColor: "#6b7280",
    color: "white",
    border: "none",
    borderRadius: "4px",
  },
  table: {
    width: "100%",
    borderCollapse: "collapse",
  },
};