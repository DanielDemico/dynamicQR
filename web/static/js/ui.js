/**
 * DynamicQR — UI & DOM Interactions
 */

const UI = {
  // Toast notifications
  showToast(message, type = 'info') {
    const container = document.getElementById('toastContainer');
    if (!container) return;

    const toast = document.createElement('div');
    toast.className = `toast toast-${type}`;

    let icon = 'ℹ️';
    if (type === 'success') icon = '✅';
    if (type === 'error') icon = '⚠️';

    toast.innerHTML = `
      <span style="font-size: 18px;">${icon}</span>
      <div style="flex: 1; word-break: break-word;">${message}</div>
    `;

    container.appendChild(toast);

    // Trigger animation
    requestAnimationFrame(() => {
      toast.classList.add('show');
    });

    setTimeout(() => {
      toast.classList.remove('show');
      setTimeout(() => toast.remove(), 350);
    }, 3500);
  },

  // Modal open/close
  openModal(modalId) {
    const modal = document.getElementById(modalId);
    if (modal) {
      modal.classList.add('active');
      const firstInput = modal.querySelector('input:not([type=hidden])');
      if (firstInput) firstInput.focus();
    }
  },

  closeModal(modalId) {
    const modal = document.getElementById(modalId);
    if (modal) {
      modal.classList.remove('active');
    }
  },

  // Copy to clipboard
  async copyToClipboard(text, successMsg = 'Link copiado para a área de transferência!') {
    try {
      if (navigator.clipboard && window.isSecureContext) {
        await navigator.clipboard.writeText(text);
      } else {
        const textArea = document.createElement('textarea');
        textArea.value = text;
        textArea.style.position = 'fixed';
        textArea.style.left = '-999999px';
        document.body.appendChild(textArea);
        textArea.focus();
        textArea.select();
        document.execCommand('copy');
        textArea.remove();
      }
      this.showToast(successMsg, 'success');
    } catch {
      this.showToast('Não foi possível copiar automaticamente.', 'error');
    }
  },

  // Date formatting
  formatDate(dateString) {
    if (!dateString) return '';
    const date = new Date(dateString);
    return date.toLocaleDateString('pt-BR', {
      day: '2-digit',
      month: '2-digit',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  },

  // Update Top Metrics
  updateMetrics(items, isHealthy) {
    const totalEl = document.getElementById('metricTotalQRs');
    const clicksEl = document.getElementById('metricTotalClicks');
    const statusEl = document.getElementById('metricServerStatus');

    if (totalEl) totalEl.textContent = items.length;

    if (clicksEl) {
      const totalClicks = items.reduce((acc, item) => acc + (item.click_count || 0), 0);
      clicksEl.textContent = totalClicks.toLocaleString('pt-BR');
    }

    if (statusEl) {
      statusEl.textContent = isHealthy ? 'Online' : 'Desconectado';
      statusEl.style.color = isHealthy ? 'var(--success)' : 'var(--danger)';
    }
  },

  // Render cards
  renderCards(items) {
    const grid = document.getElementById('qrGrid');
    if (!grid) return;

    if (!items || items.length === 0) {
      grid.innerHTML = `
        <div class="empty-state">
          <div class="empty-icon">📱</div>
          <h3>Nenhum QR Code Dinâmico</h3>
          <p>Crie seu primeiro QR Code dinâmico para poder alterar o link de destino a qualquer momento sem trocar a imagem impressa!</p>
          <button class="btn btn-primary" onclick="UI.openModal('createModal')">
            + Criar Primeiro QR Code
          </button>
        </div>
      `;
      return;
    }

    grid.innerHTML = items
      .map(
        (item) => `
      <div class="qr-card" data-id="${item.id}">
        <div class="qr-card-header">
          <div>
            <h2 class="qr-card-title">${escapeHTML(item.title)}</h2>
            <div style="margin-top: 6px; display: flex; gap: 8px;">
              <span class="badge badge-slug">/${escapeHTML(item.slug)}</span>
              <span class="badge badge-clicks">👁️ ${item.click_count || 0} scans</span>
            </div>
          </div>
        </div>

        <div class="qr-card-body">
          <div class="qr-thumbnail-wrapper" onclick="App.openViewModal('${item.id}', '${escapeHTML(item.title)}', '${item.short_url}', '${item.slug}')" title="Clique para ampliar">
            <img class="qr-thumbnail" src="${item.qr_code_image_url}" alt="QR Code ${escapeHTML(item.title)}" loading="lazy">
          </div>

          <div class="qr-card-info">
            <div class="url-box">
              <div class="url-label">Link Encurtador:</div>
              <div class="short-link-container">
                <a href="${item.short_url}" target="_blank" rel="noopener noreferrer" class="short-link-text">
                  ${item.short_url}
                </a>
                <button class="btn-copy" onclick="UI.copyToClipboard('${item.short_url}')" title="Copiar Link">
                  📋
                </button>
              </div>
            </div>

            <div class="url-box">
              <div class="url-label">Destino Atual:</div>
              <a href="${item.target_url}" target="_blank" rel="noopener noreferrer" class="target-url-link" title="${item.target_url}">
                🔗 ${truncate(item.target_url, 38)}
              </a>
            </div>
          </div>
        </div>

        <div class="qr-card-footer">
          <span class="qr-card-date">${this.formatDate(item.created_at)}</span>
          <div class="qr-card-actions">
            <button class="btn btn-secondary btn-sm" onclick="App.openEditModal('${item.id}', '${escapeHTML(item.title)}', '${escapeHTML(item.target_url)}')" title="Editar Destino">
              ✏️ Editar
            </button>
            <button class="btn btn-secondary btn-sm" onclick="App.openViewModal('${item.id}', '${escapeHTML(item.title)}', '${item.short_url}', '${item.slug}')" title="Baixar Imagem">
              ⬇️ Baixar
            </button>
            <button class="btn btn-danger btn-sm btn-icon-only" onclick="App.confirmDelete('${item.id}', '${escapeHTML(item.title)}')" title="Excluir">
              🗑️
            </button>
          </div>
        </div>
      </div>
    `
      )
      .join('');
  },
};

function escapeHTML(str) {
  if (!str) return '';
  return str.replace(/[&<>'"]/g, (tag) => ({
    '&': '&amp;',
    '<': '&lt;',
    '>': '&gt;',
    "'": '&#39;',
    '"': '&quot;',
  }[tag] || tag));
}

function truncate(str, max) {
  if (!str) return '';
  if (str.length <= max) return str;
  return str.substring(0, max) + '...';
}
