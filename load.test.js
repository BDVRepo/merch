import http from 'k6/http';
import { check, sleep } from 'k6';

const BASE_URL = 'http://localhost:8080/api';

export let options = {
  vus: 1000, 
  rps: 1000, 
  duration: '1m',
};

export default function () {
  const userId = __VU; // Уникальный номер виртуального пользователя
  const username = `loadtester${userId}`;

  // 1. Авторизация
  const token = getAuthToken(username, 'loader');
  check(token, { 'Token exists': (t) => t !== undefined });
  

  // 2. Покупка товара
  buyItem(token);

  // 3. Отправка монет следующему пользователю (циклически)
  const nextUserId = userId % 1000 + 1;
  sendCoins(token, `loadtester${nextUserId}`);

  // 4. Получение информации о себе
  getInfo(token);
}

// Функция авторизации
function getAuthToken(username,id, password) {
  const res = http.post(`${BASE_URL}/auth`, JSON.stringify({ username, password }), {
    headers: { 'Content-Type': 'application/json' },
  });
  check(res, { [`Auth success`]: (r) => r.status === 200 });

  return res.status === 200 ? JSON.parse(res.body).token : null;
}

// Функция заголовков авторизации
function getAuthHeaders(token) {
  return {
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`,
    },
  };
}

// Функция покупки товара
function buyItem(token) {
  const headers = getAuthHeaders(token);
  let buyRes = http.get(`${BASE_URL}/buy/pen`, headers);
  check(buyRes, { 'Buy success': (r) => [200, 400].includes(r.status) });
}

// Функция отправки монет
function sendCoins(token, toUsername) {
  const headers = getAuthHeaders(token);
  let sendCoinRes = http.post(
    `${BASE_URL}/sendCoin`,
    JSON.stringify({
      to_username: toUsername,
      amount: 1,
    }),
    headers
  );
  check(sendCoinRes, { 'SendCoin success': (r) => [200, 400].includes(r.status) });
}

// Функция получения информации о себе
function getInfo(token) {
  const headers = getAuthHeaders(token);
  let infoRes = http.get(`${BASE_URL}/info`, headers);
  check(infoRes, { 'Info success': (r) => [200, 401].includes(r.status) });
}
