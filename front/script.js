// script.js
const BASE_URL = "http://localhost:4444";
let jwtToken = localStorage.getItem("jwtToken") || "";
let isAuthorized = !!jwtToken;
let currentUserID = localStorage.getItem("currentUserID") || "";
let currentRoomID = null;

let ws = null; // WebSocket
let currentMovieIndex = 0;
let moviesForRoom = [];

/**
 * Универсальный запрос
 */
async function apiCall(method, url, data = null) {
    const headers = { "Content-Type": "application/json" };
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

/**
 * Сохранить токен
 */
function setToken(token) {
    jwtToken = token;
    localStorage.setItem("jwtToken", token);
    isAuthorized = true;
    renderNav();
}

/**
 * JWT -> parse user_id
 */
function parseJwt(token) {
    const base64Url = token.split('.')[1];
    const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
    const jsonPayload = decodeURIComponent(
        atob(base64)
            .split('')
            .map(c => '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2))
            .join('')
    );
    return JSON.parse(jsonPayload);
}

/**
 * WebSocket
 */
function initWebSocket() {
    if (!currentUserID) return;
    const wsUrl = `ws://${window.location.hostname}:4444/ws?user_id=${currentUserID}`;
    ws = new WebSocket(wsUrl);

    ws.onopen = () => console.log("WS connected");
    ws.onmessage = (event) => {
        const data = JSON.parse(event.data);
        console.log("WS message:", data);

        if (data.type === "INVITATION") {
            showInvitationModal(data);
        }
        else if (data.type === "MATCH") {
            // У вас мэтч
            showInstantMatchModal(data);
        }
        else if (data.type === "ROOM_ACCEPTED") {
            // Друг подтвердил комнату => у создателя showRoomSection
            console.log("ROOM_ACCEPTED from friend, room_id=", data.room_id);
            if (data.room_id) {
                currentRoomID = data.room_id;
                // Скрываем pending-section
                document.getElementById("pending-section").classList.add("hidden");
                showRoomSection();
            }
        }
    };
    ws.onerror = (err) => console.error("WS error:", err);
    ws.onclose = () => console.log("WS closed");
}

/**
 * Модалка приглашения
 */
function showInvitationModal(invData) {
    // invData = {type:"INVITATION", room_id, from_user, message?}
    const invitationModal = document.getElementById("invitation-modal");
    const invitationText = document.getElementById("invitation-text");
    invitationText.textContent = `Пользователь ${invData.from_user} пригласил вас в комнату ${invData.room_id}.`;

    invitationModal.classList.remove("hidden");

    // Принять
    const btnAccept = document.getElementById("btn-invite-accept");
    btnAccept.onclick = () => {
        invitationModal.classList.add("hidden");
        // Отправляем на сервер, что приняли
        if (ws) {
            ws.send(JSON.stringify({
                type: "ACCEPT_ROOM",
                room_id: invData.room_id
            }));
        }
        // У себя сразу start room
        currentRoomID = invData.room_id;
        showRoomSection();
    };

    // Отклонить
    const btnDecline = document.getElementById("btn-invite-decline");
    btnDecline.onclick = () => {
        invitationModal.classList.add("hidden");
        // Если нужно, можно ws.send({type:"DECLINE_ROOM", room_id:...})
    };
}

/**
 * Модалка "У вас мэтч!"
 */
function showInstantMatchModal(data) {
    // data: {type:"MATCH", imdbID, title, poster?}
    const modal = document.getElementById("instant-match-modal");
    const details = document.getElementById("instant-match-details");

    const title = data.title || data.imdbID || "Неизвестный фильм";
    details.innerHTML = `<p>Совпадение по фильму: <b>${title}</b></p>`;
    modal.classList.remove("hidden");

    document.getElementById("close-instant-match").onclick = () => {
        modal.classList.add("hidden");
    };
}

/**
 * Логин
 */
async function doLogin(username, password) {
    const resp = await apiCall("POST", BASE_URL + "/auth/login", { username, password });
    if (resp.token) {
        setToken(resp.token);
        const decoded = parseJwt(resp.token);
        if (decoded && decoded.user_id) {
            currentUserID = decoded.user_id;
            localStorage.setItem("currentUserID", currentUserID);
        }
        initWebSocket();
    } else {
        throw new Error("Токен не получен");
    }
}

/**
 * DOM Loaded
 */
document.addEventListener("DOMContentLoaded", () => {
    renderNav();
    initAuthModal();
    initCreateRoomModal();

    if (isAuthorized && currentUserID) {
        initWebSocket();
    }

    document.getElementById("skip-btn").addEventListener("click", nextMovie);
    document.getElementById("like-btn").addEventListener("click", likeMovie);
    document.getElementById("matches-btn").addEventListener("click", viewMatches);

    // Закрыть модалку "Список мэтчей"
    document.getElementById("close-matches-modal").onclick = () => {
        document.getElementById("matches-modal").classList.add("hidden");
    };
});

/**
 * Рендер меню
 */
function renderNav() {
    const navArea = document.getElementById("nav-area");
    navArea.innerHTML = "";

    if (!isAuthorized) {
        const btnLogin = document.createElement("button");
        btnLogin.classList.add("nav-btn");
        btnLogin.textContent = "Вход";
        btnLogin.onclick = () => openAuthModal("login");

        const btnReg = document.createElement("button");
        btnReg.classList.add("nav-btn");
        btnReg.textContent = "Регистрация";
        btnReg.onclick = () => openAuthModal("register");

        navArea.appendChild(btnLogin);
        navArea.appendChild(btnReg);
    } else {
        const btnLogout = document.createElement("button");
        btnLogout.classList.add("nav-btn");
        btnLogout.textContent = "Выйти";
        btnLogout.onclick = logout;

        navArea.appendChild(btnLogout);
    }
}


function openAuthModal(defaultTab = "login") {
    const authModal = document.getElementById("auth-modal");
    authModal.classList.remove("hidden");

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

function initAuthModal() {
    const authModal = document.getElementById("auth-modal");
    document.getElementById("close-auth-modal").onclick = () => {
        authModal.classList.add("hidden");
    };

    const loginTab = document.getElementById("auth-login-tab");
    const regTab = document.getElementById("auth-register-tab");
    const loginTabContent = document.getElementById("auth-login-tab-content");
    const regTabContent = document.getElementById("auth-register-tab-content");

    loginTab.onclick = () => {
        loginTab.classList.add("active");
        regTab.classList.remove("active");
        loginTabContent.classList.remove("hidden");
        regTabContent.classList.add("hidden");
    };
    regTab.onclick = () => {
        regTab.classList.add("active");
        loginTab.classList.remove("active");
        regTabContent.classList.remove("hidden");
        loginTabContent.classList.add("hidden");
    };

    document.getElementById("login-btn").onclick = async () => {
        const username = document.getElementById("login-username").value;
        const password = document.getElementById("login-password").value;
        const info = document.getElementById("login-info");
        const err = document.getElementById("login-error");
        info.textContent = "";
        err.textContent = "";
        try {
            await doLogin(username, password);
            info.textContent = "Успешный вход!";
            setTimeout(() => authModal.classList.add("hidden"), 800);
        } catch (e) {
            err.textContent = e.message;
        }
    };

    document.getElementById("register-btn").onclick = async () => {
        const regUsername = document.getElementById("reg-username").value;
        const regPassword = document.getElementById("reg-password").value;
        const info = document.getElementById("reg-info");
        const err = document.getElementById("reg-error");
        info.textContent = "";
        err.textContent = "";
        try {
            const r = await apiCall("POST", BASE_URL + "/auth/register", {
                username: regUsername,
                password: regPassword
            });
            info.textContent = `Пользователь создан (ID: ${r.user_id}). Теперь вы можете войти.`;
        } catch (ex) {
            err.textContent = ex.message;
        }
    };
}


function initCreateRoomModal() {
    const roomModal = document.getElementById("create-room-modal");
    document.getElementById("close-room-modal").onclick = () => {
        roomModal.classList.add("hidden");
    };

    const createRoomBtn = document.getElementById("create-room-btn");
    const createRoomInfo = document.getElementById("create-room-info");
    const createRoomError = document.getElementById("create-room-error");

    createRoomBtn.onclick = async () => {
        createRoomInfo.textContent = "";
        createRoomError.textContent = "";

        const genreRus = document.getElementById("select-genre").value;
        const friendID = document.getElementById("friend-id").value;

        if (!currentUserID) {
            createRoomError.textContent = "Сначала войдите!";
            return;
        }

        const genreMap = {
            "боевик": "action",
            "комедия": "comedy",
            "драма": "drama",
            "ужасы": "horror",
            "фантастика": "sci-fi"
        };
        const genreEng = genreMap[genreRus.toLowerCase()] || "action";

        try {
            const resp = await apiCall("POST", BASE_URL + "/room", {
                genre: genreEng,
                user_ids: [currentUserID, friendID]
            });
            createRoomInfo.textContent = `Комната создана (ID: ${resp.ID})`;
            currentRoomID = resp.ID;

            // Не заходим сразу в просмотр — ждём принятия
            setTimeout(() => {
                roomModal.classList.add("hidden");
                // Скрываем welcome, показываем pending
                document.getElementById("welcome-section").classList.add("hidden");
                document.getElementById("pending-section").classList.remove("hidden");
            }, 800);

        } catch (err) {
            createRoomError.textContent = err.message;
        }
    };

    const createRoomBigBtn = document.getElementById("create-room-big-btn");
    createRoomBigBtn.onclick = () => {
        if (!isAuthorized) {
            openAuthModal("login");
            return;
        }
        roomModal.classList.remove("hidden");
    };
}


async function showRoomSection() {
    document.getElementById("pending-section").classList.add("hidden");
    document.getElementById("welcome-section").classList.add("hidden");
    document.getElementById("room-section").classList.remove("hidden");

    currentMovieIndex = 0;
    moviesForRoom = [];

    try {
        const resp = await apiCall("GET", `${BASE_URL}/room/movies?room_id=${currentRoomID}&page=1`);
        moviesForRoom = resp;
    } catch (err) {
        alert("Ошибка при загрузке фильмов: " + err.message);
        return;
    }
    renderMovie();
}


function renderMovie() {
    const cont = document.getElementById("current-movie");
    cont.innerHTML = "";

    if (!moviesForRoom || moviesForRoom.length === 0) {
        cont.innerHTML = "<p>Нет фильмов</p>";
        return;
    }
    if (currentMovieIndex >= moviesForRoom.length) {
        cont.innerHTML = "<p>Фильмы закончились</p>";
        return;
    }

    const movie = moviesForRoom[currentMovieIndex];
    let title = movie.Title || "No Title";
    let year = movie.Year || "N/A";
    let poster = "https://via.placeholder.com/400x600?text=No+Poster";
    if (movie.Poster && movie.Poster.startsWith("http")) {
        poster = movie.Poster;
    }

    cont.innerHTML = `
        <h3 style="font-size:1.8rem;">${title} (${year})</h3>
        <img 
            src="${poster}" 
            alt="${title}" 
            style="max-width:400px; display:block; margin: 10px auto;"
        />
        <p style="color:gray;">${movie.imdbID}</p>
    `;
}


function nextMovie() {
    currentMovieIndex++;
    renderMovie();
}


async function likeMovie() {
    if (!currentRoomID) {
        alert("Сначала создайте комнату");
        return;
    }
    if (currentMovieIndex >= moviesForRoom.length) {
        alert("Фильмы закончились");
        return;
    }

    const movie = moviesForRoom[currentMovieIndex];
    try {
        const resp = await apiCall("POST", BASE_URL + "/room/like", {
            room_id: currentRoomID,
            imdb_id: movie.imdbID
        });
        if (resp.match) {
            alert(`У вас мэтч по фильму: ${movie.Title || movie.imdbID}`);
        } else {
            alert(`Фильм "${movie.Title}" залайкан!`);
        }
        nextMovie();
    } catch (err) {
        alert("Ошибка при лайке: " + err.message);
    }
}


async function viewMatches() {
    if (!currentRoomID) {
        alert("Комната не выбрана");
        return;
    }
    try {
        const movies = await apiCall("GET", `${BASE_URL}/room/matches?room_id=${currentRoomID}&detailed=true`);
        if (!Array.isArray(movies) || movies.length === 0) {
            alert("Пока нет мэтчей");
            return;
        }
        showMatchesModal(movies);
    } catch (err) {
        alert("Ошибка при получении мэтчей: " + err.message);
    }
}


function showMatchesModal(movies) {
    const modal = document.getElementById("matches-modal");
    const listEl = document.getElementById("matches-list");
    listEl.innerHTML = "";

    movies.forEach((m) => {
        const title = m.Title || "No Title";
        const year = m.Year || "N/A";
        let poster = "https://via.placeholder.com/200x300?text=No+Poster";
        if (m.Poster && m.Poster.startsWith("http")) {
            poster = m.Poster;
        }
        const div = document.createElement("div");
        div.style.marginBottom = "30px";
        div.innerHTML = `
            <h4>${title} (${year})</h4>
            <img src="${poster}" alt="${title}" />
            <p style="color:gray;">${m.imdbID}</p>
        `;
        listEl.appendChild(div);
    });

    modal.classList.remove("hidden");
}


function logout() {
    jwtToken = "";
    isAuthorized = false;
    currentUserID = "";
    currentRoomID = null;
    localStorage.removeItem("jwtToken");
    localStorage.removeItem("currentUserID");

    if (ws) {
        ws.close();
        ws = null;
    }
    document.getElementById("pending-section").classList.add("hidden");
    document.getElementById("room-section").classList.add("hidden");
    document.getElementById("welcome-section").classList.remove("hidden");
    renderNav();
}