// Minimal app utilities: api wrappers, modal, nav toggle, token helper
window.apiGet = async (path) => apiCall('GET', path);
window.apiPost = async (path, body) => apiCall('POST', path, body);
window.apiDelete = async (path) => apiCall('DELETE', path);

async function apiCall(method, path, body) {
  const token = localStorage.getItem('auth_token');
  const headers = { 'Accept': 'application/json' };
  if (token) headers['Authorization'] = 'Bearer ' + token;
  if (body && !(body instanceof FormData)) {
    headers['Content-Type'] = 'application/json';
    body = JSON.stringify(body);
  }
  const res = await fetch(path, { method, headers, body });
  if (!res.ok) {
    const t = await res.text();
    throw new Error(`HTTP ${res.status}: ${t}`);
  }
  // handle empty body
  const txt = await res.text();
  return txt ? JSON.parse(txt) : {};
}

// modal helpers
function showModal(content, title) {
  const root = document.getElementById('modalRoot');
  root.setAttribute('aria-hidden', 'false');
  root.innerHTML = `<div class="modal" role="dialog" aria-modal="true">
     <div style="display:flex;justify-content:space-between;align-items:center;margin-bottom:8px">
       <strong>${title||''}</strong>
       <button onclick="closeModal()" class="btn-ghost">✕</button>
     </div>
     <div>${(typeof content === 'string')? content : ''}</div>
  </div>`;
  if (content instanceof Node) {
    const modal = root.querySelector('.modal');
    modal.appendChild(content);
  }
  // wire cancel buttons in forms
  root.querySelectorAll('[data-action="cancel"]').forEach(b => b.addEventListener('click', closeModal));
}

function closeModal() {
  const root = document.getElementById('modalRoot');
  root.setAttribute('aria-hidden', 'true');
  root.innerHTML = '';
}

// nav toggle
document.addEventListener('DOMContentLoaded', () => {
  const hamb = document.getElementById('hamburger');
  const side = document.getElementById('sideNav');
  if (hamb && side) hamb.addEventListener('click', () => side.style.display = side.style.display === 'block' ? 'none' : 'block');
  // basic theme toggle
  const t = document.getElementById('toggleTheme');
  if (t) t.addEventListener('click', () => document.documentElement.classList.toggle('dark'));
});
