const express = require("express");
const app = express();
const bp = require("body-parser");
const { main } = require("./dish.js");

app.use(bp.json());
app.use(bp.urlencoded({ extended: true }));

app.get("/dish/:q", async (req, res) => {
  const data = req.params.q;
  const dish = await main(data);
  res.send(dish);
});

const PORT = 8082;
app.listen(PORT, () => {
  console.log("listening on port 8082");
});
