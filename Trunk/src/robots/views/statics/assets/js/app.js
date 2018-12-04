$(function() {
    // 读取body data-type 判断是哪个页面然后执行相应页面方法，方法在下面。
    var dataType = $('body').attr('data-type');
    console.log(dataType);
    for (key in pageData) {
        if (key == dataType) {
            pageData[key]();
        }
    }
    //     // 判断用户是否已有自己选择的模板风格
    //    if(storageLoad('SelcetColor')){
    //      $('body').attr('class',storageLoad('SelcetColor').Color)
    //    }else{
    //        storageSave(saveSelectColor);
    //        $('body').attr('class','theme-black')
    //    }

    autoLeftNav();
    $(window).resize(function() {
        autoLeftNav();
        console.log($(window).width());
    });

    //    if(storageLoad('SelcetColor')){

    //     }else{
    //       storageSave(saveSelectColor);
    //     }

    var si = loadConfig("serverInfo");
    if (si) {
        $('#serverip').val(si.serverip);
        $('#serverport').val(si.serverport);
        $('#serverid').val(si.serverid);
    }

    var ri = loadConfig("robotInfo");
    if (ri) {
        $('#account_prefix').val(ri.account_prefix);
        $('#account_start').val(ri.account_start);
        $('#password').val(ri.password);
        $('#name_prefix').val(ri.name_prefix);
        $('#robot_count').val(ri.robot_count);
    }

    jQuery.ajax({
        type: "post",
        url: "/robot/control",
        dataType: "json",
        data: '{"cmd":"queryinfo","args":""}',
        success: function(result) {
            if (result.status == 200) {
                var info = JSON.parse(result.data);
                if (info.started) {
                    updateTime = setInterval(updateInfo, 1000);
                    $('#btn-start').text("停止");
                }
            }
        }
    });
})

var echartsA;

// 页面数据
var pageData = {
    // ===============================================
    // 首页
    // ===============================================
    'index': function indexData() {
        $('#example-r').DataTable({

            bInfo: false, //页脚信息
            dom: 'ti'
        });


        // ==========================
        // 百度图表A http://echarts.baidu.com/
        // ==========================

        echartsA = echarts.init(document.getElementById('tpl-echarts'));
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
            tooltip: {
                trigger: 'axis'
            },
            grid: {
                top: '3%',
                left: '3%',
                right: '4%',
                bottom: '3%',
                containLabel: true
            },
            xAxis: [{
                type: 'category',
                boundaryGap: false,
                data: xaxis,
            }],
            yAxis: [{
                type: 'value'
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

        echartsA.setOption(option);
        echartsA.hideLoading();
    }
}

function cpuRate(ratio) {
    var options = echartsA.getOption();
    options.series[0].data.shift();
    options.series[0].data.push(ratio);
    echartsA.hideLoading();
    echartsA.setOption(options);
}
// 风格切换

$('.tpl-skiner-toggle').on('click', function() {
    $('.tpl-skiner').toggleClass('active');
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

$.fn.serializeObject = function() {
    var o = {};
    var a = this.serializeArray();
    $.each(a, function() {
        if (o[this.name] !== undefined) {
            if (!o[this.name].push) {
                o[this.name] = [o[this.name]];
            }
            o[this.name].push(this.value || '');
        } else {
            o[this.name] = this.value || '';
        }
    });
    return o;
}

var isStart = false; //启动状态
var updateTime;
var errcount = 0;

//设置服务器信息
$('#server-submit').on('click', setServerInfo);
//设置机器人信息
$('#robot-submit').on('click', setRobotInfo);
//启动或者停止机器人
$('#btn-start').on('click', function() {
    if ($('#serverip').val() == '' || $('#serverport').val() == '' || $('#serverid').val() == '') {
        msgbox("提示", "服务器信息不能为空");
        return;
    }
    if ($('#account_prefix').val() == '' ||
        $('#account_start').val() == '' ||
        $('#password').val() == '' ||
        $('#name_prefix').val() == '' ||
        $('#robot_count').val() == '') {
        msgbox("提示", "机器人信息不能为空");
        return;
    }

    if (!isStart) {
        setServerInfo();
        setRobotInfo();
        startRobot();
    } else {
        stopRobot();
    }
})

function startRobot() {
    jQuery.ajax({
        type: "post",
        url: "/robot/control",
        dataType: "json",
        data: '{"cmd":"startall","args":""}',
        success: function(result) {
            if (result.status == 200) {
                isStart = true;
                errcount = 0;
                $("#err-table").empty();
                updateTime = setInterval(updateInfo, 1000);
                $('#btn-start').text("停止");
            }
        }
    });
}

function stopRobot() {
    clearInterval(updateTime);
    jQuery.ajax({
        type: "post",
        url: "/robot/control",
        dataType: "json",
        data: '{"cmd":"stopall","args":""}',
        success: function(result) {
            if (result.status == 200) {
                isStart = false;
                $('#btn-start').text("启动");
                $('#robot-total').text($('#robot_count').val());
                $('#robot-connected').text(0);
                $('#robot-ready').text(0);
                $('#robot-errors').text(0);
            }
        }
    });
}

function sendCommand(cmd, args) {
    jQuery.ajax({
        type: "post",
        url: "/robot/control",
        dataType: "json",
        data: '{"cmd":"' + cmd + '","args":"' + args + '"}',
        success: function(result) {
            $("#op-table").append('<tr class="gradeX"><td>' + cmd + '</td><td>' + args + '</td><td>' + JSON.stringify(result) + '</td></tr>');
        }
    });
}

function updateInfo() {
    jQuery.ajax({
        type: "post",
        url: "/robot/control",
        dataType: "json",
        data: '{"cmd":"queryinfo","args":""}',
        success: function(result) {
            if (result.status == 200) {
                var info = JSON.parse(result.data);
                $('#robot-total').text(info.total);
                $('#robot-connected').text(info.connected);
                $('#robot-ready').text(info.ready);
                $('#robot-errors').text(info.errors);
                parseErrs(JSON.parse(result.errinfo));
            }
        }
    });
    jQuery.ajax({
        type: "get",
        url: "http://" + $('#serverip').val() + ":8011/sysinfo",
        dataType: "json",
        success: function(result) {
            if (result.status == 200) {
                parseSysInfo(JSON.parse(result.sysinfo));
            }
        }
    });
}

function setLoad(info, pgs, num) {
    info.text(num + '% / 100%');
    setProgress(pgs, num);
}

function setProgress(pgs, num) {
    if (num > 100) {
        num = 100;
    }
    pgs.width(num + '%')
    if (num > 90) {
        pgs.addClass('am-progress-bar-danger');
    } else if (num > 80) {
        pgs.addClass('am-progress-bar-warning');
    } else {
        pgs.removeClass('am-progress-bar-warning');
        pgs.removeClass('am-progress-bar-danger');
    }
}

function bytesToSize(bytes) {
    if (bytes === 0) return '0 B';
    var k = 1024,
        sizes = ['B', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'],
        i = Math.floor(Math.log(bytes) / Math.log(k));

    return (bytes / Math.pow(k, i)).toPrecision(3) + ' ' + sizes[i];
}
var up = 10 * 1024 * 1024; // 上行：10M
var down = 10 * 1024 * 1024; // 下行：10M

function parseSysInfo(sysinfo) {
    if (sysinfo) {
        cpuRate(sysinfo.cpuload);
        setLoad($('#cpuinfo'), $('#pgs-cpuinfo'), sysinfo.cpuload);
        setLoad($('#meminfo'), $('#pgs-meminfo'), sysinfo.memload);

        var sendload = sysinfo.netsend / down * 100;
        var recvload = sysinfo.netrecv / up * 100;

        setProgress($('#pgs-netinfo'), sendload);
        setProgress($('#pgs-netinfo-recv'), recvload);
        $('#netinfo').text("send:" + bytesToSize(sysinfo.netsend) + "/recv:" + bytesToSize(sysinfo.netrecv));
    }
}

function parseErrs(errinfo) {
    if (errinfo) {
        var err = errinfo.error;
        var len = err.length;
        if (errcount != len) {
            $("#err-table").empty();
            for (var i = 0; i < len; i++) {
                appendErr(err[i].err, err[i].account, err[i].role, err[i].index);
            }
            errcount = len;
        }
    }

}

function setServerInfo() {
    if ($('#serverip').val() == '' || $('#serverport').val() == '' || $('#serverid').val() == '') {
        msgbox("提示", "服务器信息不能为空");
        return;
    }

    var data = $("#serverinfo").serializeObject();
    jQuery.ajax({
        type: "post",
        url: "/setting/server",
        dataType: "json",
        data: JSON.stringify(data),
        success: function(result) {
            if (result.status == 200) {
                saveConfig("serverInfo", data)
            }
        }
    });
}

function setRobotInfo() {
    if ($('#account_prefix').val() == '' ||
        $('#account_start').val() == '' ||
        $('#password').val() == '' ||
        $('#name_prefix').val() == '' ||
        $('#robot_count').val() == '') {
        msgbox("提示", "机器人信息不能为空");
        return;
    }

    var data = $("#robotinfo").serializeObject();

    jQuery.ajax({
        type: "post",
        url: "/setting/robot",
        dataType: "json",
        data: JSON.stringify(data),
        success: function(result) {
            if (result.status == 200) {
                saveConfig("robotInfo", data)
            }
        }
    });
}

function saveConfig(name, data) {
    localStorage.setItem(name, JSON.stringify(data));
}

function loadConfig(name) {
    if (localStorage.getItem(name)) {
        return JSON.parse(localStorage.getItem(name));
    } else {
        return false;
    }
}

function shutdown() {
    jQuery.ajax({
        type: "post",
        url: "/robot/control",
        dataType: "json",
        data: '{"cmd":"shutdown","args":""}',
        success: function(result) {
            if (result.status == 200) {
                window.opener = null;
                window.open('', '_self');
                window.close();
            }
        }
    });
}


function appendErr(err, acc, role, index) {
    //$("#err-table").empty()
    $("#err-table").append('<tr class="gradeX"><td>' + err + '</td><td>' + acc + '</td><td>' + role + '</td><td>' + index + '</td></tr>');
}

function msgbox(title, msg) {
    $('#confirm-msg-title').text(title);
    $("#alert-box-msg").text(msg);
    $('#alert-box').modal('open');
}

$('#cmd-switch').on('click', function() {
    $('#my-prompt_title').text("切场景");
    $('#my-prompt_info').text("请输入场景号");
    $('.am-modal-prompt-input').val("");
    var $prompt = $('#my-prompt');
    var $confirmBtn = $prompt.find('#my-prompt-confirm');
    var $cancelBtn = $prompt.find('#my-prompt-cancel');
    $confirmBtn.off('click').on('click', function() {
        var scene = $('.am-modal-prompt-input').val();
        if (scene != '') {
            sendCommand("switch_scene", scene);
        } else {
            msgbox("提示", "拜托，请输入要移动的点!");
        }

    });

    $cancelBtn.off('click').on('click', function() {
        msgbox("提示", "逗我玩呢!");
    });

    $prompt.modal();
});

$('#cmd-guild').on('click', function() {
    $('#my-prompt_title').text("加入公会");
    $('#my-prompt_info').text("请输入公会编号");
    $('.am-modal-prompt-input').val("");
    var $prompt = $('#my-prompt');
    var $confirmBtn = $prompt.find('#my-prompt-confirm');
    var $cancelBtn = $prompt.find('#my-prompt-cancel');
    $confirmBtn.off('click').on('click', function() {
        var guild = $('.am-modal-prompt-input').val();
        if (guild != '') {
            sendCommand("join_guild", guild);
        } else {
            msgbox("提示", "拜托，请输入公会编号!");
        }

    });

    $cancelBtn.off('click').on('click', function() {
        msgbox("提示", "逗我玩呢!");
    });

    $prompt.modal();
});

$('#cmd-custom').on('click', function() {
    $('#my-prompt_title').text("发送自定义消息");
    $('#my-prompt_info').text("格式:消息号 [i/i64/f/d/s/o 值]");
    $('.am-modal-prompt-input').val("");
    var $prompt = $('#my-prompt');
    var $confirmBtn = $prompt.find('#my-prompt-confirm');
    var $cancelBtn = $prompt.find('#my-prompt-cancel');
    $confirmBtn.off('click').on('click', function() {
        var custom = $('.am-modal-prompt-input').val();
        if (custom != '') {
            sendCommand("send_custom", custom);
        } else {
            msgbox("提示", "拜托，请输入自定义消息!");
        }

    });

    $cancelBtn.off('click').on('click', function() {
        msgbox("提示", "逗我玩呢!");
    });

    $prompt.modal();
});

$('#cmd-gm').on('click', function() {
    $('#my-prompt_title').text("发送GM命令");
    $('#my-prompt_info').text("请输入GM命令");
    $('.am-modal-prompt-input').val("");
    var $prompt = $('#my-prompt');
    var $confirmBtn = $prompt.find('#my-prompt-confirm');
    var $cancelBtn = $prompt.find('#my-prompt-cancel');
    $confirmBtn.off('click').on('click', function() {
        var custom = $('.am-modal-prompt-input').val();
        if (custom != '') {
            sendCommand("send_gm", custom);
        } else {
            msgbox("提示", "拜托，请输入自定义消息!");
        }

    });

    $cancelBtn.off('click').on('click', function() {
        msgbox("提示", "逗我玩呢!");
    });

    $prompt.modal();
});


$('#cmd-moveto').on('click', function() {
    $('#my-prompt_title').text("瞬移");
    $('#my-prompt_info').text("请输入要移动的点x z用空格分隔");
    $('.am-modal-prompt-input').val("");
    var $prompt = $('#my-prompt');
    var $confirmBtn = $prompt.find('#my-prompt-confirm');
    var $cancelBtn = $prompt.find('#my-prompt-cancel');
    $confirmBtn.off('click').on('click', function() {
        var pos = $('.am-modal-prompt-input').val();
        if (pos != '') {
            sendCommand("moveto", pos);
        } else {
            msgbox("提示", "拜托，请输入要移动的点!");
        }

    });

    $cancelBtn.off('click').on('click', function() {
        msgbox("提示", "逗我玩呢!");
    });

    $prompt.modal();
});

$('#cmd-multiscene').on('click', function() {
    $('#my-prompt_title').text("组队副本");
    $('#my-prompt_info').text("请输入多人副本编号");
    $('.am-modal-prompt-input').val("");
    var $prompt = $('#my-prompt');
    var $confirmBtn = $prompt.find('#my-prompt-confirm');
    var $cancelBtn = $prompt.find('#my-prompt-cancel');
    $confirmBtn.off('click').on('click', function() {
        var sceneid = $('.am-modal-prompt-input').val();
        if (sceneid != '') {
            sendCommand("enter_multi_scene", sceneid);
        } else {
            msgbox("提示", "拜托，请输入具体副本编号!");
        }

    });

    $cancelBtn.off('click').on('click', function() {
        msgbox("提示", "逗我玩呢!");
    });

    $prompt.modal();
});

$('#cmd-scenemove').on('click', function() {
    $('#my-prompt_title').text("切换场景且移动起来");
    $('#my-prompt_info').text("副本号逗号隔开 副本号和持续时间空格隔开 1,2 10");
    $('.am-modal-prompt-input').val("");
    var $prompt = $('#my-prompt');
    var $confirmBtn = $prompt.find('#my-prompt-confirm');
    var $cancelBtn = $prompt.find('#my-prompt-cancel');
    $confirmBtn.off('click').on('click', function() {
        var custom = $('.am-modal-prompt-input').val();
        if (custom != '') {
            sendCommand("scene_move", custom);
        } else {
            msgbox("提示", "拜托，请输入点东西好么!");
        }

    });

    $cancelBtn.off('click').on('click', function() {
        msgbox("提示", "逗我玩呢!");
    });

    $prompt.modal();
});
$('#cmd-scenelist').on('click', function() {
    $('#my-prompt-select_title').text("进入副本");
    $('#my-prompt-select_info').text("请选择副本编号");
    $('.am-modal-prompt-select-input').val("");
    var $prompt = $('#my-prompt-select');
    var $confirmBtn = $prompt.find('#my-prompt-select-confirm');
    var $cancelBtn = $prompt.find('#my-prompt-select-cancel');
    jQuery.ajax({
        type: "post",
        url: "/robot/control",
        dataType: "json",
        async:false,
        data: '{"cmd":"get_clone_sceneid_list","args":""}',
        success: function(result) {            
            if (result.status == 200) {
                var list = result.data;
                //根据id查找对象， 
                var obj=document.getElementById('js-selected'); 
                for (let i = 0; i < list.length; i++) {
                    //添加一个选项 
                    // obj.add(new Option("文本","值")); //这个只能在IE中有效 
                    obj.options.add(new Option(list[i],list[i])); //这个兼容IE与firefox 
                }
            
            }
        }
    });

    $confirmBtn.off('click').on('click', function() {
        var args = $('#js-selected').val()
        if (args != '') {
            sendCommand("enter_clone_scene_details", args);
        } else {
            msgbox("提示", "拜托，请输入具体副本编号!");
        }

    });

    $cancelBtn.off('click').on('click', function() {
        msgbox("提示", "逗我玩呢!");
    });
    $prompt.modal();

});

$('#btn-quit').on('click', function() {
    $('#confirm-msg').text("确定退出么?")
    $('#my-confirm').modal({
        relatedTarget: this,
        onConfirm: function(options) {
            shutdown();
        },
        // closeOnConfirm: false,
        onCancel: function() {}
    });
})