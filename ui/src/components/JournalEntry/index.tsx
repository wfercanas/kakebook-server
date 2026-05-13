import type { TJournalEntry } from "../../model/journal";

import { Typography, Stack } from "@mui/material";

interface IJournalEntry {
  entry: TJournalEntry;
}

function JournalEntry(props: IJournalEntry) {
  const { entry } = props;
  return (
    <Stack
      bgcolor="#F5F5F5"
      flexDirection="row"
      justifyContent="space-between"
      gap="8px"
      padding="4px 8px"
      borderRadius="4px"
      key={entry.entry_id}
      sx={{
        "&:hover": {
          bgcolor: "#FAFAFA",
          cursor: "pointer",
        },
      }}
    >
      <Stack>
        <Typography variant="body1" fontWeight="bold">
          {entry.description}
        </Typography>
        <Typography variant="body2">{entry.date}</Typography>
      </Stack>
      <Typography variant="body2" fontWeight="600">
        {Intl.NumberFormat("es-ES", {
          style: "currency",
          currency: "COP",
        }).format(entry.amount)}
      </Typography>
    </Stack>
  );
}

export { JournalEntry };
