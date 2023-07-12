// import axios from 'axios';

// // Add a request interceptor
// axios.interceptors.request.use(
//   function (config) {
//     // Get the access token from storage
//     const accessToken = localStorage.getItem('access_token');

//     // Attach the access token to the Authorization header
//     if (accessToken) {
//       config.headers.Authorization = `Bearer ${accessToken}`;
//     }

//     return config;
//   },
//   function (error) {
//     return Promise.reject(error);
//   }
// );

// // Make API requests with Axios
// axios.get('/api/some-endpoint')
//   .then((response) => {
//     // Handle the response
//   })
//   .catch((error) => {
//     // Handle errors
//   });

// import axios from 'axios';

// // Add a response interceptor
// axios.interceptors.response.use(
//   function (response) {
//     return response;
//   },
//   function (error) {
//     const originalRequest = error.config;

//     // Check if the error status is 401 and there is no previous token refresh attempt
//     if (error.response.status === 401 && !originalRequest._retry) {
//       originalRequest._retry = true;

//       // Send a refresh token request to obtain a new access token
//       return axios.post('/api/refresh-token', { refreshToken: localStorage.getItem('refresh_token') })
//         .then((response) => {
//           // Update the access token in storage
//           localStorage.setItem('access_token', response.data.accessToken);

//           // Update the Authorization header with the new access token
//           axios.defaults.headers.common.Authorization = `Bearer ${response.data.accessToken}`;

//           // Retry the failed original request
//           return axios(originalRequest);
//         })
//         .catch((error) => {
//           // Handle token refresh error
//           // Redirect to login page or display an error message
//         });
//     }

//     return Promise.reject(error);
//   }
// );

// function logout() {
//     // Clear the access token and any other stored data
//     localStorage.removeItem('access_token');
//     localStorage.removeItem('refresh_token');

//     // Redirect to the login page or perform any other necessary actions
//   }
