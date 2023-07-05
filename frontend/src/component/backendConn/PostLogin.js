//Fetch()를 사용하는 외부 패키지 또한 비동기(asynchronous) 처리해 주어야 함!
function PostLogin(URL, user_id, password) {
  return fetch(URL, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      user_id: user_id,
      password: password,
    }),
  }).then(function (responce) {
    return responce.json();
  });
}

export default PostLogin;
