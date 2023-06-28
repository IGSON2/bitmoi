//Fetch()를 사용하는 외부 패키지 또한 비동기(asynchronous) 처리해 주어야 함!
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
