export default function authHeader() {
  const user = JSON.parse(localStorage.getItem('user'));
  const token = localStorage.getItem("token");
  if (user && token) {
    return { Authorization: token };
  } else {
    return {};
  }
}