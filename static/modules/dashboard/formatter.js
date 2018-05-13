function rankFormatter( val, row, idx ) {
    var btn_class = 'btn-default';
    if (val == 1) {
        btn_class = 'blue';
    }

    return '<button class="btn ' + btn_class + ' btn-xs">' + val + '</button>';
}

function dashboardOrgGroupNameOfStartupEventFormatter( val, row, idx ) {
    return dashboardOrgGroupNameFormatter( val, row, StartupEvent );
}

function dashboardOrgGroupNameOfShockEventFormatter( val, row, idx ) {
    return dashboardOrgGroupNameFormatter( val, row, ShockEvent );
}

function dashboardOrgGroupNameOfSpeedingEventFormatter( val, row, idx ) {
    return dashboardOrgGroupNameFormatter( val, row, SpeedingEvent );
}

function dashboardOrgGroupNameOfProximityEventFormatter( val, row, idx ) {
    return dashboardOrgGroupNameFormatter( val, row, ProximityEvent );
}

function dashboardOrgGroupNameFormatter( val, row, eventType ) {
    var groupName = row.group_name;
    if ( row.item.endsWith( "/0" ) ) {
        groupName = '<span class="font-grey-silver">Default</span>';
    }

    var asset = row.item.split( "/", 2 ),
        param = {
            stats_mode:  true,
            fast_paging: "on",
            org_id:      asset[0],
            group_id:    asset[1],
            start_date:  null,
            end_date:    null,
            event_type:  eventType
        };


    if ( $("#start_date").val() !== undefined && ! $("#start_date").prop( "disabled" ) ) {
        param.start_date = $("#start_date").val();

        // 검색 시작날짜가 설정된 경우
        if ( $("#end_date").val() !== undefined && ! $("#end_date").prop( "disabled" ) ) { // 시작/종료 날짜가 설정된 경우(기간 검색인 경우)
            param.end_date = $("#end_date").val();
        } else { // 시작 날짜만 설졍된 경우 (특정 날짜를 검색하는 경우)

        }
    }
    var prefix = '<a href="#" class="btn-show-ipaslog-on-modal" style="color: inherit; " data-query="'+ $.param( param ) + '">',
        suffix = '</a>';
    return prefix + row.org_name + '<i class="fa fa-angle-right mlr10"></i>' + groupName + suffix;
}



function dashboardIpasEquipIdOfStartupEventFormatter( val, row, idx ) {
    return dashboardIpasEquipIdFormatter( val, row, StartupEvent );
}

function dashboardIpasEquipIdOfShockEventFormatter( val, row, idx ) {
    return dashboardIpasEquipIdFormatter( val, row, ShockEvent );
}

function dashboardIpasEquipIdOfSpeedingEventFormatter( val, row, idx ) {
    return dashboardIpasEquipIdFormatter( val, row, SpeedingEvent );
}

function dashboardIpasEquipIdOfProximityEventFormatter( val, row, idx ) {
    return dashboardIpasEquipIdFormatter( val, row, ProximityEvent );
}

dashboardIpasEquipIdOfStartupEventFormatter

function dashboardIpasEquipIdFormatter( equipId, row, eventType ) {
    var param = {
        stats_mode:  true,
        fast_paging: "on",
        equip_id:    equipId,
        start_date:  null,
        end_date:    null,
        event_type:  eventType
    };
    // console.log(param);

    if ( $("#start_date").val() !== undefined && ! $("#start_date").prop( "disabled" ) ) {
        param.start_date = $("#start_date").val();

        // 검색 시작날짜가 설정된 경우
        if ( $("#end_date").val() !== undefined && ! $("#end_date").prop( "disabled" ) ) { // 시작/종료 날짜가 설정된 경우(기간 검색인 경우)
            param.end_date = $("#end_date").val();
        } else { // 시작 날짜만 설졍된 경우 (특정 날짜를 검색하는 경우)

        }
    }

    var prefix = '<a href="#" class="btn-show-ipaslog-on-modal" style="color: inherit; " data-query="'+ $.param( param ) + '">',
        suffix = '</a>';
    return prefix + getIpasTag( equipId ) + suffix;
}