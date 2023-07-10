function PostOrderJson(jsonURL, orderObject) {
  return fetch(jsonURL, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(orderObject),
  }).then(function (responce) {
    return responce.json();
  });
}

export default PostOrderJson;
