package sqlite

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

var testSeriesList = []string{
	"dc1.testhost_test_domain.cpu.0.cpu.idle",
	"dc1.testhost_test_domain.cpu.0.cpu.interrupt",
	"dc1.testhost_test_domain.cpu.0.cpu.nice",
	"dc1.testhost_test_domain.cpu.0.cpu.softirq",
	"dc1.testhost_test_domain.cpu.0.cpu.steal",
	"dc1.testhost_test_domain.cpu.0.cpu.system",
	"dc1.testhost_test_domain.cpu.0.cpu.user",
	"dc1.testhost_test_domain.cpu.0.cpu.wait",
	"dc1.testhost_test_domain.df.boot.df_complex.free",
	"dc1.testhost_test_domain.df.boot.df_complex.reserved",
	"dc1.testhost_test_domain.df.boot.df_complex.used",
	"dc1.testhost_test_domain.df.root.df_complex.free",
	"dc1.testhost_test_domain.df.root.df_complex.reserved",
	"dc1.testhost_test_domain.df.root.df_complex.used",
	"dc1.testhost_test_domain.disk.dm-0.disk_io_time.io_time",
	"dc1.testhost_test_domain.disk.dm-0.disk_io_time.weighted_io_time",
	"dc1.testhost_test_domain.disk.dm-0.disk_merged.read",
	"dc1.testhost_test_domain.disk.dm-0.disk_merged.write",
	"dc1.testhost_test_domain.disk.dm-0.disk_octets.read",
	"dc1.testhost_test_domain.disk.dm-0.disk_octets.write",
	"dc1.testhost_test_domain.disk.dm-0.disk_ops.read",
	"dc1.testhost_test_domain.disk.dm-0.disk_ops.write",
	"dc1.testhost_test_domain.disk.dm-0.disk_time.read",
	"dc1.testhost_test_domain.disk.dm-0.disk_time.write",
	"dc1.testhost_test_domain.disk.dm-0.pending_operations",
	"dc1.testhost_test_domain.entropy.entropy",
	"dc1.testhost_test_domain.hddtemp.temperature.sda",
	"dc1.testhost_test_domain.interface.eth0.if_errors.rx",
	"dc1.testhost_test_domain.interface.eth0.if_errors.tx",
	"dc1.testhost_test_domain.interface.eth0.if_octets.rx",
	"dc1.testhost_test_domain.interface.eth0.if_octets.tx",
	"dc1.testhost_test_domain.interface.eth0.if_packets.rx",
	"dc1.testhost_test_domain.interface.eth0.if_packets.tx",
	"dc1.testhost_test_domain.processes.fork_rate",
	"dc1.testhost_test_domain.processes.ps_state.blocked",
	"dc1.testhost_test_domain.processes.ps_state.paging",
	"dc1.testhost_test_domain.processes.ps_state.running",
	"dc1.testhost_test_domain.processes.ps_state.sleeping",
	"dc1.testhost_test_domain.processes.ps_state.stopped",
	"dc1.testhost_test_domain.processes.ps_state.zombies",
	"dc1.testhost_test_domain.swap.swap.cached",
	"dc1.testhost_test_domain.swap.swap.free",
	"dc1.testhost_test_domain.swap.swap.used",
	"dc1.testhost_test_domain.swap.swap_io.in",
	"dc1.testhost_test_domain.swap.swap_io.out",
}

func TestSQLite(t *testing.T) {
	sql, err := NewSQLite(`t-data/sqlite`)
	Convey("Create DB", t, func() {
		So(err, ShouldEqual, nil)
	})
	err = sql.SaveSeriesList(testSeriesList)
	Convey("SaveSeriesList", t, func() {
		So(err, ShouldEqual, nil)
	})
	err = sql.SaveFields(`dc1.testserv1`)
	Convey("SaveFieldsList", t, func() {
		So(nil, ShouldEqual, nil)
	})
}
