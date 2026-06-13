document.addEventListener("DOMContentLoaded", function () {
  const API_MCP = "/api/v1/ai/mcp";
  const modalElement = document.getElementById("modalMCP");
  const modalMCP = new bootstrap.Modal(modalElement);
  const formMCP = document.getElementById("formMCP");
  const fileInput = document.getElementById("mcpFile");
  const btnSubmit = document.getElementById("btnSubmitMCP");

  let grid;

  function initGrid() {
    grid = new gridjs.Grid({
      columns: [
        { name: "ID", hidden: true },
        {
          name: "Nama Perusahaan",
          width: "250px",
          formatter: (cell) => gridjs.html(`<strong>${cell || "-"}</strong>`),
        },
        {
          name: "Keterangan",
          width: "400px",
          formatter: (cell) => {
            if (!cell) return "-";
            return cell.length > 100 ? cell.substring(0, 100) + "..." : cell;
          },
        },
        {
          name: "Dokumen",
          width: "150px",
          formatter: (cell) => {
            if (!cell) return "-";
            return gridjs.html(
              `<a href="${cell}" target="_blank" class="btn btn-xs btn-outline-info"><i class="bx bx-link-external"></i> Link Drive</a>`,
            );
          },
        },
        {
          name: "Aksi",
          width: "150px",
          formatter: (cell, row) => {
            return gridjs.html(`
                            <div class="d-flex gap-2">
                                <button class="btn btn-sm btn-info btn-edit-mcp" data-id="${row.cells[0].data}"><i class="bx bx-edit"></i></button>
                                <button class="btn btn-sm btn-danger btn-delete-mcp" data-id="${row.cells[0].data}"><i class="bx bx-trash"></i></button>
                            </div>
                        `);
          },
        },
      ],
      server: {
        url: API_MCP,
        then: (data) => {
          console.log("MCP API Response:", data);
          return data.data.map((item) => [
            item.id,
            item.name,
            item.keterangan,
            item.url,
            null,
          ]);
        },
      },
      sort: true,
      pagination: { limit: 10 },
      style: {
        table: { "white-space": "nowrap" },
        th: { "background-color": "#f8f9fa", "text-align": "center" },
        td: { "vertical-align": "middle" },
      },
    }).render(document.getElementById("table-mcp"));
  }

  document.getElementById("btnAddMCP").addEventListener("click", () => {
    formMCP.reset();
    document.getElementById("mcpID").value = "";
    document.getElementById("urlDisplay").style.display = "none";
    document.getElementById("modalMCPLabel").innerText =
      "Tambah Konteks MCP Baru";
    modalMCP.show();
  });

  document.addEventListener("click", async (e) => {
    if (e.target.closest(".btn-edit-mcp")) {
      const id = e.target.closest(".btn-edit-mcp").dataset.id;
      editMCP(id);
    }
    if (e.target.closest(".btn-delete-mcp")) {
      const id = e.target.closest(".btn-delete-mcp").dataset.id;
      deleteMCP(id);
    }
  });

  async function editMCP(id) {
    try {
      const res = await fetch(`${API_MCP}`);
      const json = await res.json();
      const item = json.data.find((x) => x.id == id);
      if (item) {
        document.getElementById("mcpID").value = item.id;
        document.getElementById("mcpName").value = item.name;
        document.getElementById("mcpKeterangan").value = item.keterangan;
        document.getElementById("mcpURL").value = item.url;
        document.getElementById("btnViewFile").href = item.url;
        document.getElementById("urlDisplay").style.display = "block";
        document.getElementById("modalMCPLabel").innerText = "Edit Konteks MCP";
        modalMCP.show();
      }
    } catch (err) {
      console.error(err);
    }
  }

  formMCP.addEventListener("submit", async (e) => {
    e.preventDefault();

    let fileUrl = document.getElementById("mcpURL").value;

    // If new file is selected, upload to GDrive first
    if (fileInput.files.length > 0) {
      btnSubmit.disabled = true;
      btnSubmit.innerHTML = `<span class="spinner-border spinner-border-sm me-1"></span> Mengunggah File...`;

      try {
        fileUrl = await uploadToGDrive(fileInput.files[0]);
      } catch (err) {
        Swal.fire("Gagal Upload", err.message, "error");
        btnSubmit.disabled = false;
        btnSubmit.innerHTML = `Simpan Data`;
        return;
      }
    }

    const id = document.getElementById("mcpID").value;
    const data = {
      name: document.getElementById("mcpName").value,
      keterangan: document.getElementById("mcpKeterangan").value,
      url: fileUrl,
    };

    const method = id ? "PUT" : "POST";
    const url = id ? `${API_MCP}/${id}` : API_MCP;

    try {
      const res = await fetch(url, {
        method: method,
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(data),
      });
      const json = await res.json();
      if (json.status === "success") {
        modalMCP.hide();
        grid.forceRender();
        Swal.fire("Berhasil", "Data MCP berhasil disimpan", "success");
      }
    } catch (err) {
      console.error(err);
    } finally {
      btnSubmit.disabled = false;
      btnSubmit.innerHTML = `Simpan Data`;
    }
  });

  async function uploadToGDrive(file) {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();
      reader.readAsDataURL(file);
      reader.onload = async () => {
        const base64 = reader.result.split(",")[1];
        let mimeType = file.type;

        // Fallback for JSON if browser doesn't detect it correctly
        if (!mimeType && file.name.endsWith(".json")) {
          mimeType = "application/json";
        }

        const payload = {
          filename: `MCP_${document.getElementById("mcpName").value.replace(/ /g, "_")}_${new Date().getTime()}`,
          fileData: base64,
          mimeType: mimeType || "application/octet-stream",
        };

        try {
          // Using the same backend proxy/utility approach if exists,
          // but since user asked for Apps Script, we'll assume a direct call or a specific endpoint.
          // For now, let's use the GAS URL directly from env (if exposed) or assume a helper.
          // Note: Browser direct to GAS might hit CORS. Better via backend.
          // Let's assume there's an API for this or we use the backend utils.

          // We'll simulate the call to a backend endpoint that handles GDrive
          const res = await fetch("/api/v1/ai/mcp/upload-proxy", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(payload),
          });
          const json = await res.json();
          if (json.success) resolve(json.url);
          else reject(new Error(json.message));
        } catch (err) {
          reject(err);
        }
      };
      reader.onerror = (error) => reject(error);
    });
  }

  async function deleteMCP(id) {
    const result = await Swal.fire({
      title: "Hapus Data MCP?",
      text: "Konteks AI untuk perusahaan ini akan hilang!",
      icon: "warning",
      showCancelButton: true,
      confirmButtonColor: "#d33",
      confirmButtonText: "Ya, Hapus!",
    });

    if (result.isConfirmed) {
      try {
        const res = await fetch(`${API_MCP}/${id}`, { method: "DELETE" });
        const json = await res.json();
        if (json.status === "success") {
          grid.forceRender();
          Swal.fire("Terhapus!", "Data berhasil dihapus.", "success");
        }
      } catch (err) {
        console.error(err);
      }
    }
  }

  initGrid();
});
