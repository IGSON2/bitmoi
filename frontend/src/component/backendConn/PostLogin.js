function PostLogin(URL, user_id, password) {
  return fetch(URL, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      user_id: user_id,
      password: password,
    }),
  });
}

export default PostLogin;
