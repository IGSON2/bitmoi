//Fetch()를 사용하는 외부 패키지 또한 비동기(asynchronous) 처리해 주어야 함!
function PostOrderJson(jsonURL, orderObject) {
  var resultChart = [
    {
      close: 0,
      high: 0,
      low: 0,
      open: 0,
      pairname: "",
      time: 0,
      volume: 0,
    },
  ];
  return fetch(jsonURL, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(orderObject),
  })
    .then(function (responce) {
      return responce.json();
    })
    .then(function (data) {
      resultChart = data;
      return resultChart;
    });
}

export default PostOrderJson;
