import { useEffect } from "react";

function App() {
  useEffect(() => {
    fetch("/api/users/a635ba0b-4125-4e0d-93cc-378b5cd41d90")
      .then((result) => {
        if (!result.ok) {
          throw new Error("User not found");
        }
        return result.json();
      })
      .then((data) => console.log(data))
      .catch((err) => console.error(err));
  }, []);

  return <h1>Hello World</h1>;
}

export default App;
