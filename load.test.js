import http from 'k6/http';
import { check, sleep } from 'k6';

const BASE_URL = 'http://localhost:8080/api';

export let options = {
  vus: 500,  // Количество виртуальных пользователей
  rps: 1000,  // Запросов в секунду
  duration: '1m',  // Длительность теста
};

export default function () {
  // 1. Авторизация
  const token = getAuthToken('loadtester1', 'loader');
  check(token, { 'Token exists': (t) => t !== undefined });
  // 2. Покупка товара
  buyItem(token);

  // 3. Отправка монет пользователю 2
  sendCoins(token, 'loadtester2');

  // 4. Получение информации о себе
  getInfo(token);

  // Пауза между действиями (чтобы имитировать поведение пользователя)
  sleep(1);
}

// Функция авторизации
function getAuthToken(username, password) {
  const res = http.post(`${BASE_URL}/auth`, JSON.stringify({ username, password }), {
    headers: { 'Content-Type': 'application/json' },
  });
  check(res, { [`Auth ${username} success`]: (r) => r.status === 200 });

  return JSON.parse(res.body).token;
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
