/*************************************************
 * Глобальные переменные и настройки
 *************************************************/
const BASE_URL = "http://localhost:4444"; // ваш бекенд

// Сохраняем JWT-токен
let jwtToken = localStorage.getItem("jwtToken") || "";
let isAuthorized = !!jwtToken; // логическое значение

// Текущая комната (roomID), если создана
let currentRoomID = null;

// Текущий фильм, список фильмов (упрощённо)
let currentMovieIndex = 0;
let moviesForRoom = [];

/*************************************************
 * Функции для API
 *************************************************/
async function apiCall(method, url, data = null) {
    const headers = {
        "Content-Type": "application/json"
    };
    if (jwtToken) {
        headers["Authorization"] = "Bearer " + jwtToken;
    }
    const options = { method, headers };
    if (data) {
        options.body = JSON.stringify(data);
    }
    const response = await fetch(url, options);
    if (!response.ok) {
        const text = await response.text();
        throw new Error(text || "Request failed");
    }
    return response.json();
}

function setToken(token) {
    jwtToken = token;
    localStorage.setItem("jwtToken", token);
    isAuthorized = true;
    renderNav(); // перерисуем верхнее меню
}

/*************************************************
 * Инициализация приложения
 *************************************************/
document.addEventListener("DOMContentLoaded", () => {
    // Рендерим верхнее меню
    renderNav();

    // Инициализируем модалки
    initAuthModal();
    initCreateRoomModal();

    // Инициализируем кнопки "Пропустить" / "Лайк" (внутри комнаты)
    const skipBtn = document.getElementById("skip-btn");
    const likeBtn = document.getElementById("like-btn");

    skipBtn.addEventListener("click", () => {
        nextMovie();
    });

    likeBtn.addEventListener("click", () => {
        likeMovie();
    });
});

/*************************************************
 * Логика отрисовки верхнего меню
 *************************************************/
function renderNav() {
    const navArea = document.getElementById("nav-area");
    navArea.innerHTML = "";

    if (!isAuthorized) {
        // Не авторизован
        const loginBtn = document.createElement("button");
        loginBtn.classList.add("nav-btn");
        loginBtn.textContent = "Вход";
        loginBtn.onclick = () => openAuthModal("login");

        const regBtn = document.createElement("button");
        regBtn.classList.add("nav-btn");
        regBtn.textContent = "Регистрация";
        regBtn.onclick = () => openAuthModal("register");

        navArea.appendChild(loginBtn);
        navArea.appendChild(regBtn);
    } else {
        // Авторизован
        const createRoomBtn = document.createElement("button");
        createRoomBtn.classList.add("nav-btn");
        createRoomBtn.textContent = "Создать комнату";
        createRoomBtn.onclick = () => openCreateRoomModal();

        const matchesBtn = document.createElement("button");
        matchesBtn.classList.add("nav-btn");
        matchesBtn.textContent = "Посмотреть мэтчи";
        matchesBtn.onclick = () => viewMatches();

        const logoutBtn = document.createElement("button");
        logoutBtn.classList.add("nav-btn");
        logoutBtn.textContent = "Выйти";
        logoutBtn.onclick = () => logout();

        navArea.appendChild(createRoomBtn);
        navArea.appendChild(matchesBtn);
        navArea.appendChild(logoutBtn);
    }
}

/*************************************************
 * Вход / Регистрация (модалка)
 *************************************************/
function initAuthModal() {
    const authModal = document.getElementById("auth-modal");
    const closeAuthModalBtn = document.getElementById("close-auth-modal");

    const authLoginTab = document.getElementById("auth-login-tab");
    const authRegisterTab = document.getElementById("auth-register-tab");
    const loginTabContent = document.getElementById("auth-login-tab-content");
    const registerTabContent = document.getElementById("auth-register-tab-content");

    closeAuthModalBtn.onclick = () => { authModal.classList.add("hidden"); };

    // Переключение вкладок
    authLoginTab.onclick = () => {
        authLoginTab.classList.add("active");
        authRegisterTab.classList.remove("active");
        loginTabContent.classList.remove("hidden");
        registerTabContent.classList.add("hidden");
    };

    authRegisterTab.onclick = () => {
        authRegisterTab.classList.add("active");
        authLoginTab.classList.remove("active");
        registerTabContent.classList.remove("hidden");
        loginTabContent.classList.add("hidden");
    };

    // Обработчик входа
    const loginBtn = document.getElementById("login-btn");
    loginBtn.onclick = async () => {
        const loginUsername = document.getElementById("login-username").value;
        const loginPassword = document.getElementById("login-password").value;
        const loginInfo = document.getElementById("login-info");
        const loginError = document.getElementById("login-error");
        loginInfo.textContent = "";
        loginError.textContent = "";

        try {
            const resp = await apiCall("POST", BASE_URL + "/auth/login", {
                username: loginUsername,
                password: loginPassword
            });
            if (resp.token) {
                setToken(resp.token);
                loginInfo.textContent = "Успешный вход!";
                // Закрываем модалку
                setTimeout(() => {
                    authModal.classList.add("hidden");
                }, 800);
            } else {
                loginError.textContent = "Токен не получен";
            }
        } catch (err) {
            loginError.textContent = err.message;
        }
    };

    // Обработчик регистрации
    const registerBtn = document.getElementById("register-btn");
    registerBtn.onclick = async () => {
        const regUsername = document.getElementById("reg-username").value;
        const regPassword = document.getElementById("reg-password").value;
        const regInfo = document.getElementById("reg-info");
        const regError = document.getElementById("reg-error");
        regInfo.textContent = "";
        regError.textContent = "";

        try {
            const resp = await apiCall("POST", BASE_URL + "/auth/register", {
                username: regUsername,
                password: regPassword
            });
            regInfo.textContent = `Пользователь создан (ID: ${resp.user_id}). Теперь вы можете войти.`;
        } catch (err) {
            regError.textContent = err.message;
        }
    };
}

function openAuthModal(defaultTab = "login") {
    const authModal = document.getElementById("auth-modal");
    authModal.classList.remove("hidden");

    // Установим активную вкладку
    document.getElementById("auth-login-tab").classList.remove("active");
    document.getElementById("auth-register-tab").classList.remove("active");
    document.getElementById("auth-login-tab-content").classList.add("hidden");
    document.getElementById("auth-register-tab-content").classList.add("hidden");

    if (defaultTab === "login") {
        document.getElementById("auth-login-tab").classList.add("active");
        document.getElementById("auth-login-tab-content").classList.remove("hidden");
    } else {
        document.getElementById("auth-register-tab").classList.add("active");
        document.getElementById("auth-register-tab-content").classList.remove("hidden");
    }
}

/*************************************************
 * Создать комнату (модалка)
 *************************************************/
function initCreateRoomModal() {
    const roomModal = document.getElementById("create-room-modal");
    const closeRoomModalBtn = document.getElementById("close-room-modal");
    closeRoomModalBtn.onclick = () => { roomModal.classList.add("hidden"); };

    const createRoomInfo = document.getElementById("create-room-info");
    const createRoomError = document.getElementById("create-room-error");
    const createRoomBtn = document.getElementById("create-room-btn");
    createRoomBtn.onclick = async () => {
        const genre = document.getElementById("genre-select").value; // "comedy"
        const friendID = document.getElementById("friend-id").value;
        const currentUserID = localStorage.getItem("currentUserID"); // сохраненный при логине

        try {
            const resp = await apiCall("POST", BASE_URL + "/room", {
                genre: genre,
                user_ids: [ currentUserID, friendID ]
            });
            console.log("Комната создана:", resp);
        } catch (err) {
            console.error("Ошибка при создании комнаты:", err);
        }
    };

    createRoomBtn.onclick = async () => {
        createRoomInfo.textContent = "";
        createRoomError.textContent = "";

        const genreRus = document.getElementById("select-genre").value;
        const friendID = document.getElementById("friend-id").value;

        const genreMap = {
            "боевик": "action",
            "комедия": "comedy",
            "драма": "drama",
            "ужасы": "horror",
            "фантастика": "sci-fi"
        };
        const genreEng = genreMap[genreRus.toLowerCase()] || "action"; // по умолчанию

        try {
            // Запрос на создание комнаты
            const resp = await apiCall("POST", BASE_URL + "/room", {
                genre: genreEng,
                user_ids: [friendID] // упрощённо
            });
            createRoomInfo.textContent = `Комната создана (ID: ${resp._id})`;
            currentRoomID = resp._id;
            // Закрываем модалку
            setTimeout(() => {
                roomModal.classList.add("hidden");
                showRoomSection();
            }, 800);
        } catch (err) {
            createRoomError.textContent = err.message;
        }
    };
}

function openCreateRoomModal() {
    if (!isAuthorized) {
        openAuthModal("login");
        return;
    }
    document.getElementById("create-room-modal").classList.remove("hidden");
}

/*************************************************
 * Показ/скрытие секции с фильмами
 *************************************************/
async function showRoomSection() {
    // Загружаем список фильмов.
    // Допустим, берем сразу 20 штук, или по 1 (упрощённо).
    // У вас есть эндпоинт GET /room/movies?room_id=xxx
    // Но вы хотите показывать "по одному"? Мы можем просто взять все и показывать по очереди.
    moviesForRoom = [];
    currentMovieIndex = 0;

    try {
        const resp = await apiCall("GET", `${BASE_URL}/room/movies?room_id=${currentRoomID}&page=1`);
        moviesForRoom = resp; // массив фильмов
    } catch (err) {
        alert("Ошибка при загрузке фильмов: " + err.message);
        return;
    }

    // Показываем секцию комнаты
    const roomSection = document.getElementById("room-section");
    roomSection.classList.remove("hidden");

    renderMovie();
}

/*************************************************
 * Рендер одного фильма
 *************************************************/
function renderMovie() {
    const movieContainer = document.getElementById("current-movie");
    movieContainer.innerHTML = "";

    if (!moviesForRoom || moviesForRoom.length === 0) {
        movieContainer.innerHTML = "<p>Нет фильмов</p>";
        return;
    }

    if (currentMovieIndex >= moviesForRoom.length) {
        movieContainer.innerHTML = "<p>Фильмы закончились</p>";
        return;
    }

    const movie = moviesForRoom[currentMovieIndex];
    let poster = "https://via.placeholder.com/300x400?text=No+Poster";
    if (movie.Poster && movie.Poster.startsWith("http")) {
        poster = movie.Poster;
    }

    movieContainer.innerHTML = `
    <h3>${movie.Title} (${movie.Year || 'N/A'})</h3>
    <img src="${poster}" alt="${movie.Title}" style="max-width:200px; display:block; margin: 10px auto;" />
    <p style="color:gray;">${movie.imdbID}</p>
  `;
}

/*************************************************
 * Пропустить фильм
 *************************************************/
function nextMovie() {
    currentMovieIndex++;
    renderMovie();
}

/*************************************************
 * Лайк фильма
 *************************************************/
async function likeMovie() {
    if (!currentRoomID) {
        alert("Комната не создана");
        return;
    }
    if (currentMovieIndex >= moviesForRoom.length) {
        alert("Фильмы закончились");
        return;
    }

    const movie = moviesForRoom[currentMovieIndex];
    // POST /room/like c body { room_id, imdb_id }
    try {
        const resp = await apiCall("POST", BASE_URL + "/room/like", {
            room_id: currentRoomID,
            imdb_id: movie.imdbID
        });
        // Если бэкенд возвращает, например, { message: "liked", match: true/false }
        if (resp.match) {
            alert(`Мэтч по фильму: ${movie.Title}`);
        } else {
            alert(`Фильм "${movie.Title}" залайкан!`);
        }
        nextMovie();
    } catch (err) {
        alert("Ошибка при лайке: " + err.message);
    }
}

/*************************************************
 * Посмотреть мэтчи
 *************************************************/
async function viewMatches() {
    if (!currentRoomID) {
        alert("Сначала создайте комнату");
        return;
    }
    try {
        const resp = await apiCall("GET", `${BASE_URL}/room/matches?room_id=${currentRoomID}`);
        if (!Array.isArray(resp) || resp.length === 0) {
            alert("Пока нет мэтчей");
        } else {
            alert("Мэтчи (imdbID): " + resp.join(", "));
        }
    } catch (err) {
        alert("Ошибка при получении мэтчей: " + err.message);
    }
}

/*************************************************
 * Выход
 *************************************************/
function logout() {
    jwtToken = "";
    localStorage.removeItem("jwtToken");
    isAuthorized = false;
    currentRoomID = null;
    // Спрячем секцию комнаты
    document.getElementById("room-section").classList.add("hidden");
    renderNav();
}