import { Stack, Typography } from "@mui/material";
import { useEffect, useState } from "react";
import { Journal } from "./components/Journal";

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
    <Stack padding={"64px 32px"} gap="36px">
      <Typography variant="h5" fontWeight="bold">
        Kakebook
      </Typography>
      {entries.length == 0 ? "loading" : <Journal entries={entries} />}
    </Stack>
  );
}

export default App;
