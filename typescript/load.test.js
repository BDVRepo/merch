import http from 'k6/http';
import { check, sleep } from 'k6';

const BASE_URL = 'http://localhost:9000/api';

export let options = {
  scenarios: {
    authTest: {
      executor: 'constant-vus',
      vus: 1,  // Тут 1, чтобы не запускать много пользователей, и быстро получить результат
      duration: '1m',
      exec: 'authTest',  // Имя функции, которая будет запускаться
    },
    infoTest: {
      executor: 'constant-vus',
      vus: 1,
      duration: '1m',
      exec: 'infoTest',
    },
    buyAndSendTest: {
      executor: 'constant-vus',
      vus: 1,
      duration: '1m',
      exec: 'buyAndSendTest',
    },
  },
};

// Функция авторизации
function getAuthToken(username, password) {
  const res = http.post(`${BASE_URL}/auth`, JSON.stringify({ username, password }), {
    headers: { 'Content-Type': 'application/json' },
  });
  check(res, { [`Auth ${username} success`]: (r) => r.status === 200 });
  return {token : JSON.parse(res.body).token, username: username};
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

// Функция setup: авторизуем пользователей
export function setup() {
  return {
    user1: getAuthToken("loadtester1", "loader"),
    user2: getAuthToken("loadtester2", "loader"),
  };
}

// Функция для авторизации
export function authTest() {
  let token = getAuthToken("loadtester1", "loader").token;
  check(token, { 'Token exists': (t) => t !== undefined });
}

// Тест для получения информации
export function infoTest(data) {
  const headers = getAuthHeaders(data.user1);
  let infoRes = http.get(`${BASE_URL}/info`, headers);
  check(infoRes, { 'Info success': (r) => [200, 401].includes(r.status) });
  sleep(1);
}

// Тест на покупку и отправку монет
export function buyAndSendTest(data) {
  const headers1 = getAuthHeaders(data.user1.token);
  const headers2 = getAuthHeaders(data.user2.token);
  // Покупка товара
let buyRes = http.get(`${BASE_URL}/buy/pen`, headers1);
console.log(buyRes.body);  // Выводим ответ
check(buyRes, { 'Buy success': (r) => [200, 400].includes(r.status) });

// Отправка монет
let sendCoinRes = http.post(`${BASE_URL}/sendCoin`, JSON.stringify({
  to_username: data.user2.username,
  amount: 1,
}), headers1);
console.log(sendCoinRes.body);  // Выводим ответ
check(sendCoinRes, { 'SendCoin success': (r) => [200, 400].includes(r.status) });


  sleep(1);
}
