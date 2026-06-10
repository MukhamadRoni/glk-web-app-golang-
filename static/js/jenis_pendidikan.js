(function () {
    const API_JP = "/api/v1/jenis-pendidikan";
    const tbodyJP = document.querySelector("#tableJP tbody");

    // Fetch Jenis Pendidikan
    async function loadJP() {
        try {
            const res = await fetch(API_JP, { credentials: "same-origin" });
            const json = await res.json();
            
            tbodyJP.innerHTML = "";

            if (json.data && json.data.length > 0) {
                json.data.forEach((m, index) => {
                    const isChecked = m.active === 'T' ? 'checked' : '';
                    tbodyJP.innerHTML += `
                        <tr>
                            <td>${index + 1}</td>
                            <td><span class="badge bg-secondary">${m.jenis_pendidikan}</span></td>
                            <td>${m.name}</td>
                            <td>
                                <div class="form-check form-switch form-switch-md mb-3" dir="ltr">
                                    <input type="checkbox" class="form-check-input switch-active-jp" data-id="${m.id}" ${isChecked}>
                                </div>
                            </td>
                            <td>
                                <button class="btn btn-sm btn-info waves-effect waves-light btn-edit-jp" data-id="${m.id}" data-code="${m.jenis_pendidikan}" data-name="${m.name}" title="Edit">
                                    <i class="bx bx-pencil"></i>
                                </button>
                                <button class="btn btn-sm btn-danger waves-effect waves-light btn-delete-jp" data-id="${m.id}" title="Hapus">
                                    <i class="bx bx-trash"></i>
                                </button>
                            </td>
                        </tr>
                    `;
                });
            } else {
                tbodyJP.innerHTML = `<tr><td colspan="5" class="text-center">Belum ada data Jenis Pendidikan</td></tr>`;
            }
        } catch (error) {
            console.error("Error loading jenis pendidikan:", error);
        }
    }

    // Toggle Active Status
    document.addEventListener("change", async function(e) {
        if (e.target.classList.contains("switch-active-jp")) {
            const id = e.target.getAttribute("data-id");
            const newStatus = e.target.checked ? "T" : "F";
            
            try {
                const res = await fetch(`${API_JP}/${id}/active`, {
                    method: "PATCH",
                    headers: { "Content-Type": "application/json" },
                    credentials: "same-origin",
                    body: JSON.stringify({ active: newStatus })
                });
                
                if (!res.ok) {
                    const errData = await res.json();
                    Swal.fire('Gagal', "Gagal mengubah status: " + (errData.message || errData.error || "Unknown Error"), 'error');
                    e.target.checked = !e.target.checked; // Revert visually
                }
            } catch (error) {
                console.error(error);
                Swal.fire('Error', 'Terjadi kesalahan jaringan.', 'error');
                e.target.checked = !e.target.checked; // Revert visually
            }
        }
    });

    // Event Delegation for Create/Edit/Delete Buttons
    document.addEventListener("click", async function(e) {
        // Create JP
        const btnCreateJP = e.target.closest('#btnCreateJP');
        if (btnCreateJP) {
            e.preventDefault();
            const code = document.getElementById("jpCode").value;
            const name = document.getElementById("jpName").value;
            if (!code || !name) {
                Swal.fire('Peringatan', 'Kode dan Nama wajib diisi!', 'warning');
                return;
            }

            try {
                const res = await fetch(API_JP, {
                    method: "POST",
                    headers: { "Content-Type": "application/json" },
                    credentials: "same-origin",
                    body: JSON.stringify({ jenis_pendidikan: code, name: name })
                });
                if (res.ok) {
                    document.getElementById("jpCode").value = "";
                    document.getElementById("jpName").value = "";
                    await loadJP();
                    Swal.fire('Berhasil', 'Berhasil menambahkan jenis pendidikan!', 'success');
                } else {
                    const errData = await res.json();
                    Swal.fire('Gagal', "Gagal menambahkan data: " + (errData.message || errData.error || "Unknown Error"), 'error');
                }
            } catch (error) {
                console.error(error);
                Swal.fire('Error', 'Terjadi kesalahan jaringan.', 'error');
            }
        }

        // Delete JP
        const btnDeleteJP = e.target.closest('.btn-delete-jp');
        if (btnDeleteJP) {
            const id = btnDeleteJP.getAttribute('data-id');
            const confirm = await Swal.fire({
                title: 'Yakin ingin menghapus?',
                text: "Data akan dihapus permanen.",
                icon: 'warning',
                showCancelButton: true,
                confirmButtonColor: '#d33',
                cancelButtonColor: '#3085d6',
                confirmButtonText: 'Ya, hapus!',
                cancelButtonText: 'Batal'
            });
            if (!confirm.isConfirmed) return;

            try {
                const res = await fetch(`${API_JP}/${id}`, { 
                    method: "DELETE",
                    credentials: "same-origin"
                });
                if (res.ok) {
                    await loadJP();
                    Swal.fire('Terhapus!', 'Data berhasil dihapus.', 'success');
                } else {
                    Swal.fire('Gagal', 'Gagal menghapus data', 'error');
                }
            } catch (error) {
                console.error(error);
                Swal.fire('Error', 'Terjadi kesalahan jaringan.', 'error');
            }
        }

        // Edit JP
        const btnEditJP = e.target.closest('.btn-edit-jp');
        if (btnEditJP) {
            const id = btnEditJP.getAttribute('data-id');
            const currentCode = btnEditJP.getAttribute('data-code');
            const currentName = btnEditJP.getAttribute('data-name');
            
            const { value: formValues } = await Swal.fire({
                title: 'Edit Jenis Pendidikan',
                html:
                    `<input id="swal-jpCode" class="swal2-input" placeholder="Kode" value="${currentCode}">` +
                    `<input id="swal-jpName" class="swal2-input" placeholder="Nama" value="${currentName}">`,
                focusConfirm: false,
                showCancelButton: true,
                preConfirm: () => {
                    const code = document.getElementById('swal-jpCode').value;
                    const name = document.getElementById('swal-jpName').value;
                    if (!code || !name) {
                        Swal.showValidationMessage('Semua kolom wajib diisi!');
                    }
                    return { code, name };
                }
            });

            if (formValues) {
                if (formValues.code === currentCode && formValues.name === currentName) return; // No changes

                try {
                    const res = await fetch(`${API_JP}/${id}`, {
                        method: "PUT",
                        headers: { "Content-Type": "application/json" },
                        credentials: "same-origin",
                        body: JSON.stringify({ jenis_pendidikan: formValues.code, name: formValues.name })
                    });
                    if (res.ok) {
                        await loadJP();
                        Swal.fire('Tersimpan!', 'Data berhasil diubah.', 'success');
                    } else {
                        const errData = await res.json();
                        Swal.fire('Gagal', 'Gagal mengubah data: ' + (errData.message || errData.error || "Unknown Error"), 'error');
                    }
                } catch (error) {
                    console.error(error);
                    Swal.fire('Error', 'Terjadi kesalahan jaringan.', 'error');
                }
            }
        }
    });

    // Initial load
    loadJP();
})();
