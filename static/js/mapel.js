(function () {
    const API_MAPEL = "/api/v1/mapel";
    const API_JP = "/api/v1/jenis-pendidikan";
    const tbodyMapel = document.querySelector("#tableMapel tbody");
    const selectJP = document.getElementById("mapelJenisPendidikan");

    // Fetch Jenis Pendidikan for select dropdown
    async function loadJP() {
        try {
            const res = await fetch(API_JP + "?_t=" + new Date().getTime(), { 
                credentials: "same-origin",
                cache: "no-store" 
            });
            const json = await res.json();
            
            selectJP.innerHTML = '<option value="">-- Pilih Jenis Pendidikan --</option>';
            if (json && json.data && Array.isArray(json.data)) {
                json.data.forEach(jp => {
                    const opt = document.createElement('option');
                    opt.value = jp.id;
                    opt.text = `${jp.name} (${jp.jenis_pendidikan})`;
                    selectJP.appendChild(opt);
                });
            }
        } catch (error) {
            console.error("Error loading JP:", error);
        }
    }

    // Fetch Mata Pelajaran
    async function loadMapel() {
        try {
            const res = await fetch(API_MAPEL + "?_t=" + new Date().getTime(), { 
                credentials: "same-origin",
                cache: "no-store"
            });
            const json = await res.json();
            
            tbodyMapel.innerHTML = "";

            if (json.data && json.data.length > 0) {
                json.data.forEach((m, index) => {
                    const isChecked = m.active === 'Y' ? 'checked' : '';
                    const jpName = m.jenis_pendidikan_name ? m.jenis_pendidikan_name : '<span class="text-muted fst-italic">Belum Diatur</span>';
                    const jpId = m.jenis_pendidikan_id ? m.jenis_pendidikan_id : '';
                    tbodyMapel.innerHTML += `
                        <tr>
                            <td>${index + 1}</td>
                            <td>${jpName}</td>
                            <td><span class="badge bg-secondary">${m.code}</span></td>
                            <td>${m.name}</td>
                            <td>
                                <div class="form-check form-switch form-switch-md mb-3" dir="ltr">
                                    <input type="checkbox" class="form-check-input switch-active-mapel" data-id="${m.id}" ${isChecked}>
                                </div>
                            </td>
                            <td>
                                <button class="btn btn-sm btn-info waves-effect waves-light btn-edit-mapel" data-id="${m.id}" data-jp-id="${jpId}" data-code="${m.code}" data-name="${m.name}" title="Edit">
                                    <i class="bx bx-pencil"></i>
                                </button>
                                <button class="btn btn-sm btn-danger waves-effect waves-light btn-delete-mapel" data-id="${m.id}" title="Hapus">
                                    <i class="bx bx-trash"></i>
                                </button>
                            </td>
                        </tr>
                    `;
                });
            } else {
                tbodyMapel.innerHTML = `<tr><td colspan="5" class="text-center">Belum ada data Mata Pelajaran</td></tr>`;
            }
        } catch (error) {
            console.error("Error loading mapel:", error);
        }
    }

    // Toggle Active Status
    document.addEventListener("change", async function(e) {
        if (e.target.classList.contains("switch-active-mapel")) {
            const id = e.target.getAttribute("data-id");
            const newStatus = e.target.checked ? "Y" : "F";
            
            try {
                const res = await fetch(`${API_MAPEL}/${id}/active`, {
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
        // Create Mapel
        const btnCreateMapel = e.target.closest('#btnCreateMapel');
        if (btnCreateMapel) {
            e.preventDefault();
            const jpId = document.getElementById("mapelJenisPendidikan").value;
            const code = document.getElementById("mapelCode").value;
            const name = document.getElementById("mapelName").value;
            if (!jpId || !code || !name) {
                Swal.fire('Peringatan', 'Jenis Pendidikan, Kode, dan Nama wajib diisi!', 'warning');
                return;
            }

            try {
                const res = await fetch(API_MAPEL, {
                    method: "POST",
                    headers: { "Content-Type": "application/json" },
                    credentials: "same-origin",
                    body: JSON.stringify({ jenis_pendidikan_id: parseInt(jpId), code: code, name: name })
                });
                if (res.ok) {
                    document.getElementById("mapelCode").value = "";
                    document.getElementById("mapelName").value = "";
                    await loadMapel();
                    Swal.fire('Berhasil', 'Berhasil menambahkan mata pelajaran!', 'success');
                } else {
                    const errData = await res.json();
                    Swal.fire('Gagal', "Gagal menambahkan data: " + (errData.message || errData.error || "Unknown Error"), 'error');
                }
            } catch (error) {
                console.error(error);
                Swal.fire('Error', 'Terjadi kesalahan jaringan.', 'error');
            }
        }

        // Delete Mapel
        const btnDeleteMapel = e.target.closest('.btn-delete-mapel');
        if (btnDeleteMapel) {
            const id = btnDeleteMapel.getAttribute('data-id');
            const confirm = await Swal.fire({
                title: 'Yakin ingin menghapus?',
                text: "Data mata pelajaran akan dihapus secara logika (soft-delete).",
                icon: 'warning',
                showCancelButton: true,
                confirmButtonColor: '#d33',
                cancelButtonColor: '#3085d6',
                confirmButtonText: 'Ya, hapus!',
                cancelButtonText: 'Batal'
            });
            if (!confirm.isConfirmed) return;

            try {
                const res = await fetch(`${API_MAPEL}/${id}`, { 
                    method: "DELETE",
                    credentials: "same-origin"
                });
                if (res.ok) {
                    await loadMapel();
                    Swal.fire('Terhapus!', 'Data berhasil dihapus.', 'success');
                } else {
                    Swal.fire('Gagal', 'Gagal menghapus data', 'error');
                }
            } catch (error) {
                console.error(error);
                Swal.fire('Error', 'Terjadi kesalahan jaringan.', 'error');
            }
        }

        // Edit Mapel
        const btnEditMapel = e.target.closest('.btn-edit-mapel');
        if (btnEditMapel) {
            const id = btnEditMapel.getAttribute('data-id');
            const currentJpId = btnEditMapel.getAttribute('data-jp-id');
            const currentCode = btnEditMapel.getAttribute('data-code');
            const currentName = btnEditMapel.getAttribute('data-name');
            
            const jpOptions = Array.from(selectJP.options)
                .filter(opt => opt.value !== "")
                .map(opt => `<option value="${opt.value}" ${opt.value === currentJpId ? 'selected' : ''}>${opt.text}</option>`)
                .join('');

            const { value: formValues } = await Swal.fire({
                title: 'Edit Mata Pelajaran',
                html:
                    `<select id="swal-mapelJp" class="swal2-select" style="display: flex; margin: 10px auto; width: 80%;">${jpOptions}</select>` +
                    `<input id="swal-mapelCode" class="swal2-input" placeholder="Kode" value="${currentCode}">` +
                    `<input id="swal-mapelName" class="swal2-input" placeholder="Nama" value="${currentName}">`,
                focusConfirm: false,
                showCancelButton: true,
                preConfirm: () => {
                    const jpId = document.getElementById('swal-mapelJp').value;
                    const code = document.getElementById('swal-mapelCode').value;
                    const name = document.getElementById('swal-mapelName').value;
                    if (!jpId || !code || !name) {
                        Swal.showValidationMessage('Semua kolom wajib diisi!');
                    }
                    return { jpId, code, name };
                }
            });

            if (formValues) {
                if (formValues.jpId === currentJpId && formValues.code === currentCode && formValues.name === currentName) return; // No changes

                try {
                    const res = await fetch(`${API_MAPEL}/${id}`, {
                        method: "PUT",
                        headers: { "Content-Type": "application/json" },
                        credentials: "same-origin",
                        body: JSON.stringify({ jenis_pendidikan_id: parseInt(formValues.jpId), code: formValues.code, name: formValues.name })
                    });
                    if (res.ok) {
                        await loadMapel();
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
    loadJP().then(() => loadMapel());
})();
