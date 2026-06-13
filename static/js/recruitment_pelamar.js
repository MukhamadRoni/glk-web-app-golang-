document.addEventListener("DOMContentLoaded", function () {
  const API_RECRUITMENT_PELAMAR = "/api/v1/recruitment/pelamar";
  const startDateInput = document.getElementById("startDate");
  const endDateInput = document.getElementById("endDate");
  const modalElement = document.getElementById("modalDetailPelamar");
  let modalDetail;

  // Initialize modal after ensuring bootstrap is loaded and element exists
  if (modalElement && typeof bootstrap !== "undefined") {
    modalDetail = new bootstrap.Modal(modalElement);
  }

  const biodataContent = document.getElementById("detail-biodata-content");
  const tesContent = document.getElementById("detail-tes-content");

  if (!startDateInput || !endDateInput) return; // Prevent errors if on different page

  // Default date: today
  const today = new Date().toISOString().split("T")[0];
  startDateInput.value = today;
  endDateInput.value = today;

  let grid;

  function initGrid() {
    if (typeof gridjs === "undefined") {
      console.warn("Grid.js is not loaded. Skipping table initialization.");
      return;
    }
    grid = new gridjs.Grid({
      columns: [
        { name: "ID", hidden: true },
        {
          name: "#",
          width: "50px",
          formatter: (cell) => gridjs.html(`<b>${cell}</b>`),
        },
        { name: "Nama Pelamar", width: "200px" },
        { name: "Email", width: "200px" },
        { name: "Posisi Dilamar", width: "250px" },
        {
          name: "Status",
          width: "150px",
          formatter: (cell) => {
            let badgeClass = "bg-warning";
            if (cell === "Selesai Tes") badgeClass = "bg-info";
            if (cell === "Diterima") badgeClass = "bg-success";
            if (cell === "Ditolak") badgeClass = "bg-danger";
            return gridjs.html(
              `<span class="badge ${badgeClass}">${cell}</span>`,
            );
          },
        },
        { name: "Tanggal", width: "150px" },
        {
          name: "Aksi",
          width: "100px",
          sort: false,
          formatter: (cell, row) => {
            return gridjs.html(`
                            <button class="btn btn-sm btn-primary btn-detail-pelamar" data-id="${row.cells[0].data}">
                                <i class="bx bx-show"></i> Detail
                            </button>
                        `);
          },
        },
      ],
      pagination: {
        limit: 10,
        server: {
          url: (prev, page, limit) =>
            `${prev}${prev.includes("?") ? "&" : "?"}limit=${limit}&offset=${page * limit}`,
        },
      },
      search: {
        server: {
          url: (prev, keyword) =>
            `${prev}${prev.includes("?") ? "&" : "?"}search=${keyword}`,
        },
      },
      sort: false,
      server: {
        url: `${API_RECRUITMENT_PELAMAR}?start_date=${startDateInput.value}&end_date=${endDateInput.value}`,
        then: (data) =>
          data.data.items.map((item, index) => [
            item.ID,
            (data.data.offset || 0) + index + 1,
            item.nama_lengkap || item.pelamar.name,
            item.pelamar.email,
            `${item.target_jenjang.name} - ${item.target_mapel.name}`,
            item.status,
            new Date(item.CreatedAt).toLocaleDateString("id-ID"),
            null,
          ]),
        total: (data) => data.data.total,
      },
      language: {
        search: { placeholder: "Cari nama atau email..." },
        pagination: {
          previous: "Sebelumnya",
          next: "Selanjutnya",
          showing: "Menampilkan",
          results: () => "data",
        },
      },
      style: {
        table: { "white-space": "nowrap" },
        th: {
          "background-color": "#f8f9fa",
          color: "#495057",
          "font-weight": "bold",
          "text-align": "center",
        },
        td: { "vertical-align": "middle" },
      },
    }).render(document.getElementById("table-recruitment-pelamar"));
  }

  const refreshGrid = () => {
    grid
      .updateConfig({
        server: {
          url: `${API_RECRUITMENT_PELAMAR}?start_date=${startDateInput.value}&end_date=${endDateInput.value}`,
          then: (data) =>
            data.data.items.map((item, index) => [
              item.ID,
              (data.data.offset || 0) + index + 1,
              item.nama_lengkap || item.pelamar.name,
              item.pelamar.email,
              `${item.target_jenjang.name} - ${item.target_mapel.name}`,
              item.status,
              new Date(item.CreatedAt).toLocaleDateString("id-ID"),
              null,
            ]),
          total: (data) => data.data.total,
        },
      })
      .forceRender();
  };

  startDateInput.addEventListener("change", refreshGrid);
  endDateInput.addEventListener("change", refreshGrid);

  document.addEventListener("click", async function (e) {
    const btn = e.target.closest(".btn-detail-pelamar");
    if (btn) {
      const id = btn.getAttribute("data-id");
      showDetail(id);
    }
  });

  async function showDetail(id) {
    biodataContent.innerHTML = `<div class="text-center py-4"><div class="spinner-border text-primary" role="status"></div></div>`;
    tesContent.innerHTML = `<div class="text-center py-4"><div class="spinner-border text-primary" role="status"></div></div>`;

    // Reset and Add Action Buttons to Modal Header
    const existingButtons = document.querySelectorAll(".modal-status-btn");
    existingButtons.forEach((btn) => btn.remove());

    const modalHeader = document.querySelector(
      "#modalDetailPelamar .modal-header",
    );
    const btnContainer = document.createElement("div");
    btnContainer.className = "d-flex gap-2 ms-auto me-3 modal-status-btn";
    btnContainer.innerHTML = `
        <button class="btn btn-sm btn-outline-success" onclick="updateStatus(${id}, 'Diterima')">Set Diterima</button>
        <button class="btn btn-sm btn-outline-warning" onclick="updateStatus(${id}, 'Pending')">Set Pending</button>
        <button class="btn btn-sm btn-outline-danger" onclick="updateStatus(${id}, 'Ditolak')">Set Ditolak</button>
    `;
    modalHeader.insertBefore(
      btnContainer,
      modalHeader.querySelector(".btn-close"),
    );

    if (modalDetail) modalDetail.show();

    try {
      const res = await fetch(`${API_RECRUITMENT_PELAMAR}/${id}`);
      const json = await res.json();
      if (json.success) {
        document
          .getElementById("modalDetailPelamar")
          .setAttribute("data-current-id", id);
        renderBiodata(json.data.lamaran);
        renderHasilTes(json.data.testResults, json.data.lamaran.koreksi_nilai);
      } else {
        biodataContent.innerHTML = `<div class="alert alert-danger">Gagal mengambil data: ${json.message}</div>`;
        tesContent.innerHTML = `<div class="alert alert-danger">Gagal mengambil data: ${json.message}</div>`;
      }
    } catch (error) {
      console.error(error);
      biodataContent.innerHTML = `<div class="alert alert-danger">Terjadi kesalahan jaringan.</div>`;
    }
  }

  function renderBiodata(l) {
    biodataContent.innerHTML = `
            <div class="row">
                <div class="col-md-6">
                    <table class="table table-sm table-borderless">
                        <tr><th width="150">Nama Lengkap</th><td>: ${l.nama_lengkap}</td></tr>
                        <tr><th>Email</th><td>: ${l.pelamar.email}</td></tr>
                        <tr><th>No. WA</th><td>: ${l.no_wa}</td></tr>
                        <tr><th>Jenis Kelamin</th><td>: ${l.jenis_kelamin}</td></tr>
                        <tr><th>Alamat</th><td>: ${l.alamat_domisili}</td></tr>
                        <tr><th>Kota/Kec</th><td>: ${l.kota.name} / ${l.kecamatan.name}</td></tr>
                    </table>
                </div>
                <div class="col-md-6">
                    <table class="table table-sm table-borderless">
                        <tr><th width="150">Universitas</th><td>: ${l.universitas}</td></tr>
                        <tr><th>Prodi</th><td>: ${l.program_studi}</td></tr>
                        <tr><th>Jenjang/Sem</th><td>: ${l.jenjang_ditempuh} / ${l.semester}</td></tr>
                        <tr><th>Target Posisi</th><td>: ${l.target_jenjang.name} - ${l.target_mapel.name}</td></tr>
                        <tr><th>Ketersediaan</th><td>: ${l.ketersediaan}</td></tr>
                        <tr><th>Fee Harapan</th><td>: ${l.fee_harapan}</td></tr>
                    </table>
                </div>
            </div>

            <hr>
            <div class="row">
                <div class="col-md-12">
                    <h6 class="fw-bold"><i class="bx bx-calendar-check"></i> Ketersediaan & Jadwal Mengajar</h6>

                    <div class="row mb-3">
                        <div class="col-md-6">
                            <table class="table table-sm table-bordered text-center">
                                <thead class="table-light">
                                    <tr>
                                        <th>Metode</th>
                                        <th>Status</th>
                                    </tr>
                                </thead>
                                <tbody>
                                    ${(() => {
                                      let ket = {};
                                      try {
                                        // Handle both JSON string and comma separated fallback
                                        if (
                                          l.ketersediaan &&
                                          l.ketersediaan.startsWith("{")
                                        ) {
                                          ket = JSON.parse(l.ketersediaan);
                                        } else {
                                          const parts = (
                                            l.ketersediaan || ""
                                          ).split(",");
                                          ket = {
                                            Online: parts[0] || "-",
                                            Offline: parts[1] || "-",
                                          };
                                        }
                                      } catch (e) {
                                        ket = { Online: "-", Offline: "-" };
                                      }

                                      return Object.keys(ket)
                                        .map(
                                          (key) => `
                                            <tr>
                                                <td class="text-start ps-3">${key}</td>
                                                <td>
                                                    <span class="badge ${ket[key] === "Bersedia" ? "bg-success" : "bg-danger"}">
                                                        ${ket[key]}
                                                    </span>
                                                </td>
                                            </tr>
                                        `,
                                        )
                                        .join("");
                                    })()}
                                </tbody>
                            </table>
                        </div>
                    </div>

                    <div class="table-responsive">
                        <table class="table table-sm table-bordered text-center">
                            <thead class="table-light">
                                <tr>
                                    <th>Hari</th>
                                    <th>08.00 - 12.00</th>
                                    <th>12.00 - 18.00</th>
                                    <th>18.00 - 21.00</th>
                                    <th>Jadwal Penuh</th>
                                </tr>
                            </thead>
                            <tbody>
                                ${(() => {
                                  let jadwal = {};
                                  try {
                                    jadwal = JSON.parse(l.jadwal_free || "{}");
                                  } catch (e) {}
                                  const days = [
                                    "Senin",
                                    "Selasa",
                                    "Rabu",
                                    "Kamis",
                                    "Jumat",
                                    "Sabtu",
                                    "Minggu",
                                  ];
                                  const slots = [
                                    "08.00 - 12.00",
                                    "12.00 - 18.00",
                                    "18.00 - 21.00",
                                    "Jadwal Penuh",
                                  ];

                                  return days
                                    .map((day) => {
                                      const dayJadwal = jadwal[day] || [];
                                      return `
                                            <tr>
                                                <td class="table-light fw-bold">${day}</td>
                                                ${slots
                                                  .map(
                                                    (slot) => `
                                                    <td>
                                                        ${dayJadwal.includes(slot) ? '<i class="bx bx-check-square text-success fs-4"></i>' : '<i class="bx bx-minus text-muted"></i>'}
                                                    </td>
                                                `,
                                                  )
                                                  .join("")}
                                            </tr>
                                        `;
                                    })
                                    .join("");
                                })()}
                            </tbody>
                        </table>
                    </div>
                </div>
            </div>

            <hr>
            <div class="row">
                <div class="col-12">
                    <h6>Dokumen Lampiran:</h6>
                    <div class="d-flex gap-2">
                        <a href="${l.cv_url}" target="_blank" class="btn btn-outline-info btn-sm">
                            <i class="bx bx-file"></i> Lihat CV
                        </a>
                        <a href="${l.transkrip_url}" target="_blank" class="btn btn-outline-info btn-sm">
                            <i class="bx bx-file"></i> Lihat Transkrip
                        </a>
                    </div>
                </div>
            </div>
        `;
  }

  function renderHasilTes(tr, koreksiStr) {
    if (!tr.finished) {
      tesContent.innerHTML = `<div class="alert alert-warning">Pelamar belum menyelesaikan tes akademik.</div>`;
      return;
    }

    let answers = {};
    try {
      answers =
        typeof tr.answers === "string" ? JSON.parse(tr.answers) : tr.answers;
    } catch (e) {
      console.error("Error parsing answers:", e);
    }

    let corrections = {};
    try {
      if (koreksiStr) corrections = JSON.parse(koreksiStr);
    } catch (e) {
      console.error("Error parsing corrections:", e);
    }

    let html = `<div class="d-flex justify-content-between align-items-center mb-3">
                    <h6 class="mb-0">Hasil Tes Akademik (${tr.bankSoal.title})</h6>
                    <button class="btn btn-sm btn-success" id="btnSaveCorrection">
                        <i class="bx bx-save"></i> Simpan Koreksi
                    </button>
                </div>`;

    html += `<div class="table-responsive"><table class="table table-bordered table-sm">
            <thead class="table-light">
                <tr><th width="50">#</th><th>Pertanyaan</th><th>Jawaban Pelamar</th><th>Kunci Jawaban</th><th width="100">Status</th><th width="150">Koreksi</th></tr>
            </thead>
            <tbody>`;

    let correctCount = 0;
    let totalCount = tr.bankSoal.questions ? tr.bankSoal.questions.length : 0;

    if (tr.bankSoal.questions) {
      tr.bankSoal.questions.forEach((q, i) => {
        const questionId = q.id || q.ID;
        const qIDStr = String(questionId);
        const pelamarAns = answers[qIDStr] || answers[questionId] || "-";
        let correctAns = "-";
        if (q.options) {
          q.options.forEach((opt) => {
            if (opt.is_correct === "T") correctAns = opt.option_text;
          });
        }

        let isCorrect;
        if (corrections[qIDStr]) {
          isCorrect = corrections[qIDStr] === "T";
        } else {
          isCorrect = pelamarAns !== "-" && pelamarAns === correctAns;
        }

        if (isCorrect) correctCount++;

        html += `<tr>
                    <td>${i + 1}</td>
                    <td>${q.question_text}</td>
                    <td>${pelamarAns}</td>
                    <td>${correctAns}</td>
                    <td class="text-center status-col" data-qid="${qIDStr}">
                        ${isCorrect ? '<span class="text-success"><i class="bx bx-check-circle"></i> Benar</span>' : '<span class="text-danger"><i class="bx bx-x-circle"></i> Salah</span>'}
                    </td>
                    <td>
                        <select class="form-select form-select-sm select-correction" data-qid="${qIDStr}">
                            <option value="">-- Auto --</option>
                            <option value="T" ${corrections[qIDStr] === "T" ? "selected" : ""}>Benar (Manual)</option>
                            <option value="F" ${corrections[qIDStr] === "F" ? "selected" : ""}>Salah (Manual)</option>
                        </select>
                    </td>
                </tr>`;
      });
    }

    const score =
      totalCount > 0 ? ((correctCount / totalCount) * 100).toFixed(2) : 0;
    html += `</tbody><tfoot class="table-light">
                <tr><th colspan="4" class="text-end">Total Benar:</th><th class="text-center" id="totalCorrect">${correctCount} / ${totalCount}</th><th></th></tr>
                <tr><th colspan="4" class="text-end">Skor Akhir:</th><th class="text-center text-primary" style="font-size: 1.2rem;" id="finalScore">${score}</th><th></th></tr>
            </tfoot></table></div>`;
    tesContent.innerHTML = html;

    document.querySelectorAll(".select-correction").forEach((select) => {
      select.addEventListener("change", function () {
        updateLiveScore(tr, answers);
      });
    });

    document
      .getElementById("btnSaveCorrection")
      .addEventListener("click", function () {
        saveCorrections();
      });
  }

  function updateLiveScore(tr, answers) {
    let correctCount = 0;
    const totalCount = tr.bankSoal.questions.length;

    tr.bankSoal.questions.forEach((q) => {
      const questionId = q.id || q.ID;
      const qIDStr = String(questionId);
      const select = document.querySelector(
        `.select-correction[data-qid="${qIDStr}"]`,
      );
      const statusCell = document.querySelector(
        `.status-col[data-qid="${qIDStr}"]`,
      );

      let isCorrect;
      if (select.value === "T") {
        isCorrect = true;
        statusCell.innerHTML =
          '<span class="text-success fw-bold"><i class="bx bx-check-circle"></i> Benar (M)</span>';
      } else if (select.value === "F") {
        isCorrect = false;
        statusCell.innerHTML =
          '<span class="text-danger fw-bold"><i class="bx bx-x-circle"></i> Salah (M)</span>';
      } else {
        const pelamarAns = answers[qIDStr] || answers[questionId] || "-";
        let correctAns = "-";
        q.options.forEach((opt) => {
          if (opt.is_correct === "T") correctAns = opt.option_text;
        });
        isCorrect = pelamarAns !== "-" && pelamarAns === correctAns;
        statusCell.innerHTML = isCorrect
          ? '<span class="text-success"><i class="bx bx-check-circle"></i> Benar</span>'
          : '<span class="text-danger"><i class="bx bx-x-circle"></i> Salah</span>';
      }
      if (isCorrect) correctCount++;
    });

    const score =
      totalCount > 0 ? ((correctCount / totalCount) * 100).toFixed(2) : 0;
    document.getElementById("totalCorrect").innerText =
      `${correctCount} / ${totalCount}`;
    document.getElementById("finalScore").innerText = score;
  }

  async function saveCorrections() {
    const lamaranId = document
      .getElementById("modalDetailPelamar")
      .getAttribute("data-current-id");
    const corrections = {};
    document.querySelectorAll(".select-correction").forEach((select) => {
      if (select.value) {
        corrections[select.getAttribute("data-qid")] = select.value;
      }
    });

    try {
      Swal.fire({
        title: "Menyimpan...",
        allowOutsideClick: false,
        didOpen: () => {
          Swal.showLoading();
        },
      });
      const res = await fetch(
        `${API_RECRUITMENT_PELAMAR}/${lamaranId}/correction`,
        {
          method: "PATCH",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ corrections: JSON.stringify(corrections) }),
        },
      );
      const json = await res.json();
      if (json.success) {
        Swal.fire("Berhasil", "Koreksi nilai berhasil disimpan!", "success");
      } else {
        Swal.fire("Gagal", json.message, "error");
      }
    } catch (error) {
      console.error(error);
      Swal.fire("Error", "Terjadi kesalahan jaringan.", "error");
    }
  }

  // Global Function for Status Update
  window.updateStatus = async function (id, status) {
    const confirm = await Swal.fire({
      title: `Set status ke ${status}?`,
      text:
        status === "Ditolak"
          ? "Akun pelamar akan dinonaktifkan (Soft Delete)."
          : "",
      icon: "question",
      showCancelButton: true,
      confirmButtonText: "Ya, Lanjutkan",
    });

    if (confirm.isConfirmed) {
      try {
        Swal.fire({
          title: "Memproses...",
          allowOutsideClick: false,
          didOpen: () => {
            Swal.showLoading();
          },
        });
        const res = await fetch(`${API_RECRUITMENT_PELAMAR}/${id}/status`, {
          method: "PATCH",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ status: status }),
        });
        const json = await res.json();
        if (json.success) {
          Swal.fire(
            "Berhasil",
            `Status diperbarui ke ${status}`,
            "success",
          ).then(() => {
            if (typeof grid !== "undefined" && grid) {
              grid.forceRender();
            } else {
              location.reload(); // Fallback for dashboard
            }
          });
        } else {
          Swal.fire("Gagal", json.message, "error");
        }
      } catch (error) {
        console.error(error);
        Swal.fire("Error", "Terjadi kesalahan jaringan.", "error");
      }
    }
  };

  if (document.getElementById("table-recruitment-pelamar")) {
    initGrid();
  }
});
