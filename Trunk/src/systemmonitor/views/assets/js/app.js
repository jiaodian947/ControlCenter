$(function() {
    // 读取body data-type 判断是哪个页面然后执行相应页面方法，方法在下面。
    var dataType = $('body').attr('data-type');
    console.log(dataType);
    for (key in pageData) {
        if (key == dataType) {
            pageData[key]();
        }
    }

    autoLeftNav();
    $(window).resize(function() {
        autoLeftNav();
        console.log($(window).width())
    });

    $("#btn-start").hide(); //禁用
})

var samplePids = []

var sceneCharts = {}

function confirmPid() {
    samplePids = []
    $('#pid-list tr td input:checkbox').each(function() {
        if ($(this).is(":checked")) {
            samplePids.push($(this).val());
        }
    });

    if (samplePids.length == 0) {
        $("#btn-start").hide();
    } else {
        $("#btn-start").show();
        $("#btn-start").html('<span class="am-icon-play"></span> 开始');
    }
    console.log(samplePids);
}

var queryTimer = 0;
var isStart = false;

$('#btn-start').on('click', function() {
    if (!isStart) {
        startSample();
    } else {
        stopSample();
    }
})

function startSample() {
    jQuery.ajax({
        type: "post",
        url: "/sample",
        dataType: "json",
        data: '{"pids":"' + samplePids.join(',') + '"}',
        success: function(result) {
            if (result.status == 200) {
                initChart();
                queryTimer = setInterval(querySample, 1000);
                gameInfoTimer = setInterval(queryGameInfo, 5 * 1000);
                $("#btn-start").html('<span class="am-icon-stop"></span> 停止');
                isStart = true;
            }
        }
    });
}

function stopSample() {
    clearInterval(queryTimer);
    clearInterval(gameInfoTimer);
    jQuery.ajax({
        type: "get",
        url: "/stop",
        dataType: "json",
        success: function(result) {
            if (result.status == 200) {
                $("#btn-start").html('<span class="am-icon-play"></span> 开始');
                isStart = false;
            }
        }
    });
}

function queryGameInfo() {
    jQuery.ajax({
        type: "get",
        url: "/info",
        dataType: "json",
        success: function(result) {
            if (result.status == 200) {
                $("#total-players").text(result.data.TotalUser);
                $("#hour-count").text("0");
                $("#online-players").text(result.data.Online);
                $("#online-avg").text(result.data.Average);
                $("#max-online").text(result.data.MaxOnline);
                $("#max-online-time").text(result.data.MaxTime);
                updateSceneInfo(result.data.Scenes);
            }
        }
    });
}

var plistloaded = false;

var processInfo = {};

function refreshPid() {
    plistloaded = false;
    $("#plist").html("正在加载中，请稍候");
    jQuery.ajax({
        type: "get",
        url: "/pid",
        success: function(result) {
            if (result.status == 200) {
                plistloaded = true;
                errcount = 0;
                $("#plist").empty();
                var ps = JSON.parse(result.sysinfo);
                processInfo = {};
                for (pi in ps) {
                    processInfo[ps[pi].Pid] = ps[pi];
                    var pid = ps[pi].Pid
                    var cmd = ps[pi].Cmd.substr(0, 40)
                    $("#plist").append('<tr class="gradeX"><td>' + pid + '</td><td width="60%">' + cmd + '</td><td><input id="cb_' + pid + '" type="checkbox" value="' + pid + '"></td></tr>');
                }
            }
        }
    });
}

function querySample() {
    jQuery.ajax({
        type: "get",
        url: "/query",
        dataType: "json",
        success: function(result) {
            if (result.status == 200) {
                //var info = JSON.parse(result.data);
                console.log(result.data);
                for (var p in result.data) {
                    updateCpuInfo(result.data[p].Pid, result.data[p].Usr, result.data[p].Rss);
                }
            }
        }
    });
}

function rnd(n, m) {
    var random = Math.floor(Math.random() * (m - n + 1) + n);
    return random;
}

function isEmptyObject(obj) {
    for (var key in obj) {
        return false
    };
    return true
};

var cpucharts = {};
var chartnum = 0;

function bytesToSize(bytes) {
    if (bytes === 0) return '0 B';
    var k = 1024,
        sizes = ['B', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'],
        i = Math.floor(Math.log(bytes) / Math.log(k));

    return (bytes / Math.pow(k, i)).toPrecision(3) + ' ' + sizes[i];
}

function formatMem(num) {
    var letter = 'K'

    num = num * 4
    if (num >= 1000) {
        num = (num + 512) / 1024
        letter = 'M'
        if (num >= 10000) {
            num = (num + 512) / 1024
            letter = 'G'
        }
    }
    return '' + Math.floor(num) + letter
}

function updateCpuInfo(pid, ratio, memory) {
    if (cpucharts[pid] == undefined) {
        if (createChart(pid) == undefined)
            return;
    }
    var options = cpucharts[pid].getOption();
    options.series[0].data.shift();
    options.series[0].data.push(Math.floor(ratio));
    options.title[0].text = "内存占用:" + formatMem(memory)
    cpucharts[pid].hideLoading();
    cpucharts[pid].setOption(options);
}

function initChart() {

    for (var k in cpucharts) {
        if (cpucharts[k] == undefined) {
            continue;
        }

        cpucharts[k].dispose();
    }

    cpucharts = [];
    var $monitor = $('#monitor-chart');
    $monitor.empty();

    for (var k in samplePids) {
        addChartWidget(samplePids[k], samplePids[k]);
        createChart(samplePids[k]);
    }
}

function createChart(pid) {
    var $chart = $('#chart-' + pid.toString());
    if ($chart.length == 0) {
        return undefined;
    }
    echart = echarts.init($chart.get(0));
    var xaxis = [];
    var ydata = [];
    for (var i = 0; i < 61; i++) {
        if ((i % 10) == 0) {
            xaxis.unshift(i);
        } else {
            xaxis.unshift('');
        }
        ydata.push(0);
    }
    option = {
        title: {
            text: '内存占用：0M',
            x: 'right',
            textStyle: {
                //文字颜色
                color: '#FFF',
                //字体风格,'normal','italic','oblique'
                fontStyle: 'normal',
                //字体粗细 'normal','bold','bolder','lighter',100 | 200 | 300 | 400...
                fontWeight: 'bold',
                //字体大小
                fontSize: 18
            }
        },
        tooltip: {
            trigger: 'axis'
        },
        grid: {
            top: '30px',
            left: '10px',
            right: '10px',
            bottom: '0px',
            containLabel: true
        },
        xAxis: [{
            type: 'category',
            boundaryGap: false,
            data: xaxis,
        }],
        yAxis: [{
            type: 'value',
            min: 0,
            max: 100
        }],
        textStyle: {
            color: '#838FA1'
        },
        series: [{
            name: 'CPU负载',
            type: 'line',
            stack: '总量',
            animation: false,
            showSymbol: false,
            areaStyle: { normal: {} },
            data: ydata,
            itemStyle: {
                normal: {
                    color: '#1cabdb',
                    borderColor: '#1cabdb',
                    borderWidth: '2',
                    borderType: 'solid',
                    opacity: '1'
                },
                emphasis: {

                }
            }
        }]
    };

    echart.setOption(option);
    echart.hideLoading();
    cpucharts[pid] = echart;
    return echart;
}

function addChartWidget(name, pid) {
    var info = processInfo[pid]
    var $widget =
        '<div class="am-u-sm-4 am-u-md-4">\
        <div class="widget am-cf">\
        <div class="widget-head am-cf">\
            <div class="widget-title am-fl">监控-' + info.Cmd.substr(0, 25) + '(PID:' +
        pid +
        ')</div>\
            <div class="widget-function am-fr">\
                <a href="javascript:;" class="am-icon-cog"></a>\
            </div>\
        </div>\
        <div class="widget-body-md widget-body tpl-amendment-echarts am-fr" id="chart-' + pid + '">\
        </div>\
        </div>\
    </div>';

    $('#monitor-chart').append($widget);
    //setTimeout('createChart(' + pid.toString() + ')', 1000)
}

function updateSceneInfo(sceneinfo) {
    var option = sceneCharts.getOption()
    option.xAxis[0].data = []
    option.series[0].data = []
    if (sceneinfo.length > 0) {
        for (var index in sceneinfo) {
            var scene = sceneinfo[index]
            option.xAxis[0].data.push(scene.SceneId)
                //option.series[0].data.push(scene.Players)
            option.series[0].data.push(rnd(1, 500));
            sceneCharts.hideLoading();
            sceneCharts.setOption(option)
        }
    }

}

// 页面数据
var pageData = {
    // ===============================================
    // 首页
    // ===============================================
    'index': function indexData() {
        sceneCharts = echarts.init(document.getElementById('scene-players'));
        optionC = {
            tooltip: {
                trigger: 'axis'
            },
            animation: false,
            legend: {
                data: ['在线人数']
            },
            xAxis: [{
                type: 'category',
                axisLabel: {
                    interval: 0,
                    rotate: 90
                },
                data: []
            }],
            yAxis: [{
                type: 'value',
                name: '人数',
                min: 0,
                max: 500,
                interval: 100,
                axisLabel: {
                    formatter: '{value}'
                }
            }],
            series: [{
                name: '人数',
                type: 'bar',
                data: [],
                itemStyle: {
                    normal: {
                        color: function(params) {
                            if (params.value >= 300) {
                                return '#ff0000'
                            } else if (params.value > 150) {
                                return '#f37b1d'
                            }
                            return '#1cabdb'
                        }
                    }
                }
            }]
        };

        sceneCharts.setOption(optionC);
    }
}

// 风格切换
$('.tpl-skiner-toggle').on('click', function() {
    $('.tpl-skiner').toggleClass('active');
})

$('.tpl-skiner-content-bar').find('span').on('click', function() {
    $('body').attr('class', $(this).attr('data-color'))
    saveSelectColor.Color = $(this).attr('data-color');
    // 保存选择项
    storageSave(saveSelectColor);

})

// 侧边菜单开关
function autoLeftNav() {
    $('.tpl-header-switch-button').on('click', function() {
        if ($('.left-sidebar').is('.active')) {
            if ($(window).width() > 1024) {
                $('.tpl-content-wrapper').removeClass('active');
            }
            $('.left-sidebar').removeClass('active');
        } else {

            $('.left-sidebar').addClass('active');
            if ($(window).width() > 1024) {
                $('.tpl-content-wrapper').addClass('active');
            }
        }
    })

    if ($(window).width() < 1024) {
        $('.left-sidebar').addClass('active');
    } else {
        $('.left-sidebar').removeClass('active');
    }
}


// 侧边菜单
$('.sidebar-nav-sub-title').on('click', function() {
    $(this).siblings('.sidebar-nav-sub').slideToggle(80)
        .end()
        .find('.sidebar-nav-sub-ico').toggleClass('sidebar-nav-sub-ico-rotate');
})

$('#my-popup').on('open.modal.amui', function() {
    if (!plistloaded) {
        refreshPid()
    }
});

$('#my-popup').on('close.modal.amui', function() {
    console.log("close");
});