document.addEventListener('DOMContentLoaded', function () {
    const API_CHAT = "/api/v1/ai/chat";
    const chatContainer = document.getElementById('chatContainer');
    const welcomeMsg = document.getElementById('welcomeMsg');
    const formChat = document.getElementById('formChat');
    const chatInput = document.getElementById('chatInput');
    const btnSend = document.getElementById('btnSend');
    const btnClear = document.getElementById('btnClearChat');
    const selectMCP = document.getElementById('chatMCP');
    const selectSkill = document.getElementById('chatSkill');

    const sessionId = Date.now().toString();

    formChat.addEventListener('submit', async (e) => {
        e.preventDefault();
        const message = chatInput.value.trim();
        if (!message) return;

        // 1. UI: Remove welcome, append user message
        if (welcomeMsg) welcomeMsg.style.display = 'none';
        appendMessage('user', message);
        chatInput.value = '';

        // 2. UI: Show typing indicator
        const typingId = showTyping();

        // 3. API Call
        try {
            const res = await fetch(API_CHAT, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    message: message,
                    mcp_id: parseInt(selectMCP.value) || 0,
                    skill_id: parseInt(selectSkill.value) || 0,
                    session_id: sessionId
                })
            });

            const json = await res.json();
            hideTyping(typingId);

            if (json.success) {
                appendMessage('ai', json.message);
            } else {
                appendMessage('ai', "Maaf, terjadi kesalahan: " + json.message);
            }
        } catch (err) {
            console.error(err);
            hideTyping(typingId);
            appendMessage('ai', "Maaf, koneksi ke server terputus.");
        }
    });

    btnClear.addEventListener('click', () => {
        chatContainer.innerHTML = '';
        if (welcomeMsg) {
            chatContainer.appendChild(welcomeMsg);
            welcomeMsg.style.display = 'block';
        }
    });

    function appendMessage(role, text) {
        const div = document.createElement('div');
        div.className = `message message-${role} animate__animated animate__fadeInUp`;
        div.innerHTML = `
            <div class="message-content">
                ${text.replace(/\n/g, '<br>')}
            </div>
            <div class="text-muted font-size-10 mt-1 ${role === 'user' ? 'text-end' : ''}">
                ${new Date().toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })}
            </div>
        `;
        chatContainer.appendChild(div);
        chatContainer.scrollTop = chatContainer.scrollHeight;
    }

    function showTyping() {
        const id = 'typing-' + Date.now();
        const div = document.createElement('div');
        div.id = id;
        div.className = 'message message-ai animate__animated animate__fadeInUp';
        div.innerHTML = `
            <div class="message-content">
                <span class="typing-indicator"><i class="bx bx-loader-alt bx-spin me-1"></i> Aivene sedang mengetik...</span>
            </div>
        `;
        chatContainer.appendChild(div);
        chatContainer.scrollTop = chatContainer.scrollHeight;
        return id;
    }

    function hideTyping(id) {
        const el = document.getElementById(id);
        if (el) el.remove();
    }
});
