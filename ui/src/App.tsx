import { Stack, Typography } from "@mui/material";
import { useEffect, useState } from "react";

interface Entry {
  description: string;
  date: string;
  amount: number;
  project_id: string;
  entry_id: string;
}

function App() {
  const [entries, setEntries] = useState<Entry[]>([]);

  useEffect(() => {
    fetch("/api/projects/1b879b0b-4a22-4d49-b318-c95827b697d9/journal")
      .then((result) => {
        if (!result.ok) {
          throw new Error("Journal not found");
        }
        return result.json();
      })
      .then((data) => setEntries(data))
      .catch((err) => console.error(err));
  }, []);

  return (
    <Stack padding={"64px 32px"}>
      <Typography variant="h6">Journal</Typography>
      {entries.length > 0 &&
        entries.map((entry) => (
          <Stack
            flexDirection="row"
            justifyContent="space-between"
            gap="8px"
            padding="4px 6px"
            key={entry.entry_id}
          >
            <Typography variant="body2">
              {entry.description} - {entry.date}
            </Typography>
            <Typography variant="body2" fontWeight="600">
              {Intl.NumberFormat("us-US", {
                style: "currency",
                currency: "USD",
              }).format(entry.amount)}
            </Typography>
          </Stack>
        ))}
    </Stack>
  );
}

export default App;
