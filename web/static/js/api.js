/**
 * DynamicQR — API Client
 */

const API = {
  baseUrl: '/api/v1/qr-codes',

  async list(search = '', page = 1, limit = 20) {
    const params = new URLSearchParams();
    if (search) params.append('search', search);
    if (page) params.append('page', page);
    if (limit) params.append('limit', limit);

    const res = await fetch(`${this.baseUrl}?${params.toString()}`);
    if (!res.ok) {
      const err = await res.json().catch(() => ({}));
      throw new Error(err.error || 'Erro ao carregar lista de QR Codes');
    }
    return res.json();
  },

  async create(payload) {
    const res = await fetch(this.baseUrl, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(payload),
    });

    const data = await res.json().catch(() => ({}));
    if (!res.ok) {
      throw new Error(data.error || 'Falha ao criar QR Code');
    }
    return data;
  },

  async get(id) {
    const res = await fetch(`${this.baseUrl}/${id}`);
    if (!res.ok) {
      const err = await res.json().catch(() => ({}));
      throw new Error(err.error || 'QR Code não encontrado');
    }
    return res.json();
  },

  async update(id, payload) {
    const res = await fetch(`${this.baseUrl}/${id}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(payload),
    });

    const data = await res.json().catch(() => ({}));
    if (!res.ok) {
      throw new Error(data.error || 'Falha ao atualizar QR Code');
    }
    return data;
  },

  async delete(id) {
    const res = await fetch(`${this.baseUrl}/${id}`, {
      method: 'DELETE',
    });

    if (!res.ok && res.status !== 204) {
      const err = await res.json().catch(() => ({}));
      throw new Error(err.error || 'Falha ao excluir QR Code');
    }
    return true;
  },

  async checkHealth() {
    try {
      const res = await fetch('/healthz');
      return res.ok;
    } catch {
      return false;
    }
  }
};
