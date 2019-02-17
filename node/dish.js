const axios = require("axios");

const appid = "fbf3d292";
const appkey = "65c8bf0b798fc3ee638cef53c4fe190d";
const recipe = [];
const ingredient = [];

const random_duration = () => {
  const dur = [15, 30, 45, 60, 75, 90, 105, 120];
  const pos = Math.floor(Math.random() * 9);
  return dur[pos];
};

const dish = async food_item => {
  const food = food_item.split(" ").join("+");
  const url =
    "https://api.edamam.com/search?q=" +
    food +
    "&app_id=a176c55b&app_key=e5df3c1124e09bb097d9e5db197497ee&from=0&to=5";
  try {
    const response = await axios.get(url);
    const data = response.data;
    data.hits.forEach(hit => {
      hit.recipe.ingredients.forEach(ingredient_each => {
        ingredient.push(ingredient_each);
      });
      let random_dur = random_duration();
      recipe.push({
        label: hit.recipe.label,
        servings: hit.recipe.yield,
        ingredients: ingredient,
        ingredientLines: hit.recipe.ingredientLines,
        dietLabels: hit.recipe.dietLabels,
        healthLabels: hit.recipe.healthLabels,
        image_url: hit.recipe.image,
        duration: random_dur
      });
    });
    return recipe;
  } catch (error) {
    console.log(error.message);
  }
};

const main = async food_item => {
  const data = await dish(food_item);
  return data;
};

module.exports = { main };
