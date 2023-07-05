function PostSignUp(URL, userInfo) {
  return fetch(URL, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(userInfo),
  }).then(function (responce) {
    return responce.json();
  });
}

export default PostLogin;
