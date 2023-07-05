const accessTokenKey = "access_token";
const refreshTokenKey = "refresh_token";

function SaveAccessToken(token) {
  localStorage.setItem(accessTokenKey, token);
}

function LoadAccessToken() {
  return localStorage.getItem(accessTokenKey);
}

function SaveRefreshToken(token) {
  localStorage.setItem(refreshTokenKey, token);
}

function LoadRefreshToken() {
  return localStorage.getItem(refreshTokenKey);
}

function RemoveTokens() {
  localStorage.removeItem(accessTokenKey);
  localStorage.removeItem(refreshTokenKey);
}

export {
  SaveAccessToken,
  LoadAccessToken,
  SaveRefreshToken,
  LoadRefreshToken,
  RemoveTokens,
};
