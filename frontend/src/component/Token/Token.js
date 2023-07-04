const accessTokenKey = "access_token";

function SaveAccessToken(token) {
  localStorage.setItem(accessTokenKey, token);
}

function LoadAccessToken() {
  return localStorage.getItem(accessTokenKey);
}

export { SaveAccessToken, LoadAccessToken };
