import type { TJournalEntry } from "../../model/journal";

import { Typography, Stack, Button } from "@mui/material";
import { JournalEntry } from "../JournalEntry";
import { useState } from "react";

interface IJournal {
  entries: TJournalEntry[];
}

function Journal(props: IJournal) {
  const { entries } = props;
  const increment = 10;

  const [list, setList] = useState<TJournalEntry[]>(entries.slice(0, 5));

  function allowMore() {
    return list.length + increment <= entries.length;
  }

  function allowLess() {
    return list.length > increment;
  }

  function handleMore() {
    setList(entries.slice(0, list.length + increment));
  }

  function handleLess() {
    setList(entries.slice(0, list.length - increment));
  }

  return (
    <Stack gap="24px">
      <Stack gap="4px">
        <Typography variant="h6">Journal</Typography>
        {list.length > 0 &&
          list.map((entry) => (
            <JournalEntry key={entry.entry_id} entry={entry} />
          ))}
      </Stack>
      <Stack flexDirection="row" justifyContent="flex-end" gap="8px">
        <Button onClick={handleLess} disabled={!allowLess()}>
          Show Less
        </Button>
        <Button
          onClick={handleMore}
          disabled={!allowMore()}
          variant="contained"
        >
          Show More
        </Button>
      </Stack>
    </Stack>
  );
}

export { Journal };
