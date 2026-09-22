/**
 * DynamicQR — Main Application Logic
 */

const App = {
  items: [],
  searchTimeout: null,

  async init() {
    this.bindEvents();
    await this.loadData();

    // Auto-refresh metrics every 30 seconds
    setInterval(() => this.loadData(false), 30000);
  },

  bindEvents() {
    // Search input with debounce
    const searchInput = document.getElementById('searchInput');
    if (searchInput) {
      searchInput.addEventListener('input', (e) => {
        clearTimeout(this.searchTimeout);
        this.searchTimeout = setTimeout(() => {
          this.loadData(true, e.target.value);
        }, 300);
      });
    }

    // Create Form Submit
    const createForm = document.getElementById('createForm');
    if (createForm) {
      createForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        await this.handleCreate();
      });
    }

    // Edit Form Submit
    const editForm = document.getElementById('editForm');
    if (editForm) {
      editForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        await this.handleEdit();
      });
    }

    // Close modals on backdrop click
    document.querySelectorAll('.modal-backdrop').forEach((backdrop) => {
      backdrop.addEventListener('click', (e) => {
        if (e.target === backdrop) {
          backdrop.classList.remove('active');
        }
      });
    });

    // Close on Escape key
    document.addEventListener('keydown', (e) => {
      if (e.key === 'Escape') {
        document.querySelectorAll('.modal-backdrop.active').forEach((m) => m.classList.remove('active'));
      }
    });
  },

  async loadData(showLoading = true, search = '') {
    try {
      const isHealthy = await API.checkHealth();
      const res = await API.list(search);
      this.items = res.data || [];

      UI.updateMetrics(this.items, isHealthy);
      UI.renderCards(this.items);
    } catch (err) {
      UI.showToast(err.message || 'Erro ao carregar QR codes', 'error');
    }
  },

  async handleCreate() {
    const titleInput = document.getElementById('createTitle');
    const targetInput = document.getElementById('createTarget');
    const slugInput = document.getElementById('createSlug');
    const submitBtn = document.getElementById('btnSubmitCreate');

    const title = titleInput.value.trim();
    let target = targetInput.value.trim();
    const slug = slugInput.value.trim();

    if (!title) {
      UI.showToast('Informe um título para o QR Code', 'error');
      return;
    }

    if (!target) {
      UI.showToast('Informe o link de destino', 'error');
      return;
    }

    // Auto-prefix https if missing scheme
    if (!target.startsWith('http://') && !target.startsWith('https://')) {
      target = 'https://' + target;
      targetInput.value = target;
    }

    submitBtn.disabled = true;
    submitBtn.textContent = 'Gerando QR Code...';

    try {
      const payload = {
        title,
        target_url: target,
      };
      if (slug) {
        payload.custom_slug = slug;
      }

      const created = await API.create(payload);
      UI.showToast(`QR Code "${created.title}" criado com sucesso!`, 'success');
      UI.closeModal('createModal');

      createForm.reset();
      await this.loadData();
    } catch (err) {
      UI.showToast(err.message, 'error');
    } finally {
      submitBtn.disabled = false;
      submitBtn.textContent = 'Criar QR Code';
    }
  },

  openEditModal(id, title, targetUrl) {
    document.getElementById('editId').value = id;
    document.getElementById('editTitle').value = title;
    document.getElementById('editTarget').value = targetUrl;
    UI.openModal('editModal');
  },

  async handleEdit() {
    const id = document.getElementById('editId').value;
    const title = document.getElementById('editTitle').value.trim();
    let target = document.getElementById('editTarget').value.trim();
    const submitBtn = document.getElementById('btnSubmitEdit');

    if (!title || !target) {
      UI.showToast('Preencha todos os campos obrigatórios', 'error');
      return;
    }

    if (!target.startsWith('http://') && !target.startsWith('https://')) {
      target = 'https://' + target;
    }

    submitBtn.disabled = true;
    submitBtn.textContent = 'Salvando...';

    try {
      await API.update(id, {
        title,
        target_url: target,
      });

      UI.showToast('Destino atualizado com sucesso! Novos acessos já usarão este link.', 'success');
      UI.closeModal('editModal');
      await this.loadData();
    } catch (err) {
      UI.showToast(err.message, 'error');
    } finally {
      submitBtn.disabled = false;
      submitBtn.textContent = 'Salvar Alterações';
    }
  },

  openViewModal(id, title, shortUrl, slug) {
    document.getElementById('viewModalTitle').textContent = title;
    document.getElementById('viewShortUrl').textContent = shortUrl;
    document.getElementById('viewShortUrl').href = shortUrl;

    const img = document.getElementById('viewLargeImage');
    img.src = `/api/v1/qr-codes/${id}/image?t=${Date.now()}`;

    const downloadBtn = document.getElementById('btnDownloadPng');
    downloadBtn.onclick = () => {
      window.location.href = `/api/v1/qr-codes/${id}/image?download=true`;
    };

    UI.openModal('viewModal');
  },

  async confirmDelete(id, title) {
    if (!confirm(`Deseja realmente excluir o QR Code "${title}"?\n\nApós a exclusão, qualquer leitura do QR Code exibirá página 404.`)) {
      return;
    }

    try {
      await API.delete(id);
      UI.showToast(`QR Code "${title}" removido com sucesso.`, 'info');
      await this.loadData();
    } catch (err) {
      UI.showToast(err.message, 'error');
    }
  },
};

// Start application when DOM is ready
document.addEventListener('DOMContentLoaded', () => {
  App.init();
});
