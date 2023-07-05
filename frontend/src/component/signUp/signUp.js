import { useState } from "react";

function signUp() {
  const [userId, setUserId] = useState("");
  const [password, setPassword] = useState("");
  const [repeatPassword, setRepeatPassword] = useState("");
  const [fullName, setFullName] = useState("");
  const [email, setEmaiol] = useState("");
  // "user_id"
  // "oauth_uid"
  // "full_name"
  // "hashed_password"
  // "email"
  // "photo_url"

  // res = email
  return (
    <div>
      <form>
        <input
          type="text"
          value={userId}
          placeholder="ID"
          pattern="^[a-zA-Z0-9]+$"
        />
        {/*ID 중복체크*/}
        <input
          type="password"
          value={password}
          placeholder="password"
          pattern="^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,16}$"
        />
        <input
          type="password"
          value={repeatPassword}
          placeholder="repeat password"
          pattern="^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,16}$"
        />
        {/*PW 일치*/}
      </form>
    </div>
  );
}
