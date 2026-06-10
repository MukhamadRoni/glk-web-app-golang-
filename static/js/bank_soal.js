(function () {
    const API_BANK = "/api/v1/bank-soal";
    const tbodyBank = document.querySelector("#tableBankSoal tbody");

    async function loadBank() {
        try {
            const res = await fetch(API_BANK + "?_t=" + new Date().getTime(), { 
                credentials: "same-origin",
                cache: "no-store" 
            });
            const json = await res.json();
            
            tbodyBank.innerHTML = "";

            if (json && json.data && json.data.length > 0) {
                json.data.forEach((b, index) => {
                    const isChecked = b.active === 'T' ? 'checked' : '';
                    tbodyBank.innerHTML += `
                        <tr>
                            <td>${index + 1}</td>
                            <td>${b.title}</td>
                            <td>${b.jenis_pendidikan ? b.jenis_pendidikan.name : '-'}</td>
                            <td>${b.mata_pelajaran ? b.mata_pelajaran.name : '-'}</td>
                            <td><span class="badge bg-info">v${b.version}</span></td>
                            <td>
                                <div class="form-check form-switch form-switch-md mb-3" dir="ltr">
                                    <input type="checkbox" class="form-check-input switch-active-bank" data-id="${b.id}" ${isChecked}>
                                </div>
                            </td>
                            <td>
                                <a href="/admin/recruitment/master/bank-soal/form/${b.id}" class="btn btn-sm btn-info waves-effect waves-light" title="Edit / Versi Baru">
                                    <i class="bx bx-pencil"></i>
                                </a>
                                <button class="btn btn-sm btn-danger waves-effect waves-light btn-delete-bank" data-id="${b.id}" title="Hapus">
                                    <i class="bx bx-trash"></i>
                                </button>
                            </td>
                        </tr>
                    `;
                });
            } else {
                tbodyBank.innerHTML = `<tr><td colspan="7" class="text-center">Belum ada data Bank Soal</td></tr>`;
            }
        } catch (error) {
            console.error("Error loading bank soal:", error);
        }
    }

    document.addEventListener("change", async function(e) {
        if (e.target.classList.contains("switch-active-bank")) {
            const id = e.target.getAttribute("data-id");
            const newStatus = e.target.checked ? "T" : "F";
            
            try {
                const res = await fetch(`${API_BANK}/${id}/active`, {
                    method: "PATCH",
                    headers: { "Content-Type": "application/json" },
                    credentials: "same-origin",
                    body: JSON.stringify({ active: newStatus })
                });
                
                if (!res.ok) {
                    const errData = await res.json();
                    Swal.fire('Gagal', "Gagal mengubah status: " + (errData.message || errData.error || "Unknown Error"), 'error');
                    e.target.checked = !e.target.checked;
                }
            } catch (error) {
                console.error(error);
                Swal.fire('Error', 'Terjadi kesalahan jaringan.', 'error');
                e.target.checked = !e.target.checked;
            }
        }
    });

    document.addEventListener("click", async function(e) {
        const btnDelete = e.target.closest('.btn-delete-bank');
        if (btnDelete) {
            const id = btnDelete.getAttribute('data-id');
            const confirm = await Swal.fire({
                title: 'Yakin ingin menghapus?',
                text: "Bank soal beserta semua pertanyaan dan pilihan jawaban akan dihapus.",
                icon: 'warning',
                showCancelButton: true,
                confirmButtonColor: '#d33',
                cancelButtonColor: '#3085d6',
                confirmButtonText: 'Ya, hapus!',
                cancelButtonText: 'Batal'
            });
            if (!confirm.isConfirmed) return;

            try {
                const res = await fetch(`${API_BANK}/${id}`, { 
                    method: "DELETE",
                    credentials: "same-origin"
                });
                if (res.ok) {
                    await loadBank();
                    Swal.fire('Terhapus!', 'Data berhasil dihapus.', 'success');
                } else {
                    Swal.fire('Gagal', 'Gagal menghapus data', 'error');
                }
            } catch (error) {
                console.error(error);
                Swal.fire('Error', 'Terjadi kesalahan jaringan.', 'error');
            }
        }
    });

    loadBank();
})();
