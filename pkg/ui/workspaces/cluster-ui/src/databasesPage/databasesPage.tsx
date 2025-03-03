// Copyright 2021 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

import React from "react";
import { Link, RouteComponentProps } from "react-router-dom";
import { Tooltip } from "antd";
import "antd/lib/tooltip/style";
import classNames from "classnames/bind";
import classnames from "classnames/bind";

import { Anchor } from "src/anchor";
import { StackIcon } from "src/icon/stackIcon";
import { Pagination, ResultsPerPageLabel } from "src/pagination";
import { BooleanSetting } from "src/settings/booleanSetting";
import {
  ColumnDescriptor,
  ISortedTablePagination,
  SortedTable,
  SortSetting,
} from "src/sortedtable";
import * as format from "src/util/format";

import styles from "./databasesPage.module.scss";
import sortableTableStyles from "src/sortedtable/sortedtable.module.scss";
import {
  baseHeadingClasses,
  statisticsClasses,
} from "src/transactionsPage/transactionsPageClasses";
import { syncHistory, tableStatsClusterSetting } from "src/util";
import booleanSettingStyles from "../settings/booleanSetting.module.scss";
import { CircleFilled } from "../icon";
import LoadingError from "../sqlActivity/errorComponent";
import { Loading } from "../loading";

const cx = classNames.bind(styles);
const sortableTableCx = classNames.bind(sortableTableStyles);
const booleanSettingCx = classnames.bind(booleanSettingStyles);

// We break out separate interfaces for some of the nested objects in our data
// both so that they can be available as SortedTable rows and for making
// (typed) test assertions on narrower slices of the data.
//
// The loading and loaded flags help us know when to dispatch the appropriate
// refresh actions.
//
// The overall structure is:
//
//   interface DatabasesPageData {
//     loading: boolean;
//     loaded: boolean;
//     lastError: Error;
//     sortSetting: SortSetting;
//     databases: { // DatabasesPageDataDatabase[]
//       loading: boolean;
//       loaded: boolean;
//       name: string;
//       sizeInBytes: number;
//       tableCount: number;
//       rangeCount: number;
//       nodesByRegionString: string;
//       missingTables: { // DatabasesPageDataMissingTable[]
//         loading: boolean;
//         name: string;
//       }[];
//     }[];
//   }
export interface DatabasesPageData {
  loading: boolean;
  loaded: boolean;
  lastError: Error;
  databases: DatabasesPageDataDatabase[];
  sortSetting: SortSetting;
  automaticStatsCollectionEnabled?: boolean;
  showNodeRegionsColumn?: boolean;
}

export interface DatabasesPageDataDatabase {
  loading: boolean;
  loaded: boolean;
  lastError: Error;
  name: string;
  sizeInBytes: number;
  tableCount: number;
  rangeCount: number;
  missingTables: DatabasesPageDataMissingTable[];
  // String of nodes grouped by region in alphabetical order, e.g.
  // regionA(n1,n2), regionB(n3)
  nodesByRegionString?: string;
  numIndexRecommendations: number;
}

// A "missing" table is one for which we were unable to gather size and range
// count information during the backend call to fetch DatabaseDetails. We
// expose it here so that the component has the opportunity to try to refresh
// those properties for the table directly.
export interface DatabasesPageDataMissingTable {
  // Note that we don't need a "loaded" property here because we expect
  // the reducers supplying our properties to remove a missing table from
  // the list once we've loaded its data.
  loading: boolean;
  name: string;
}

export interface DatabasesPageActions {
  refreshDatabases: () => void;
  refreshDatabaseDetails: (database: string) => void;
  refreshTableStats: (database: string, table: string) => void;
  refreshSettings: () => void;
  refreshNodes?: () => void;
  onSortingChange?: (
    name: string,
    columnTitle: string,
    ascending: boolean,
  ) => void;
}

export type DatabasesPageProps = DatabasesPageData &
  DatabasesPageActions &
  RouteComponentProps<unknown>;

interface DatabasesPageState {
  pagination: ISortedTablePagination;
  lastDetailsError: Error;
}

class DatabasesSortedTable extends SortedTable<DatabasesPageDataDatabase> {}

export class DatabasesPage extends React.Component<
  DatabasesPageProps,
  DatabasesPageState
> {
  constructor(props: DatabasesPageProps) {
    super(props);

    this.state = {
      pagination: {
        current: 1,
        pageSize: 20,
      },
      lastDetailsError: null,
    };

    const { history } = this.props;
    const searchParams = new URLSearchParams(history.location.search);
    const ascending = (searchParams.get("ascending") || undefined) === "true";
    const columnTitle = searchParams.get("columnTitle") || undefined;
    const sortSetting = this.props.sortSetting;

    if (
      this.props.onSortingChange &&
      columnTitle &&
      (sortSetting.columnTitle != columnTitle ||
        sortSetting.ascending != ascending)
    ) {
      this.props.onSortingChange("Databases", columnTitle, ascending);
    }
  }

  componentDidMount(): void {
    this.refresh();
  }

  componentDidUpdate(): void {
    this.refresh();
  }

  private refresh(): void {
    if (this.props.refreshNodes != null) {
      this.props.refreshNodes();
    }

    if (this.props.refreshSettings != null) {
      this.props.refreshSettings();
    }

    if (
      !this.props.loaded &&
      !this.props.loading &&
      this.props.lastError === undefined
    ) {
      return this.props.refreshDatabases();
    }

    let lastDetailsError: Error;
    this.props.databases.forEach(database => {
      if (database.lastError !== undefined) {
        lastDetailsError = database.lastError;
      }
      if (
        lastDetailsError &&
        this.state.lastDetailsError?.name != lastDetailsError?.name
      ) {
        this.setState({ lastDetailsError: lastDetailsError });
      }
      if (
        !database.loaded &&
        !database.loading &&
        database.lastError === undefined
      ) {
        return this.props.refreshDatabaseDetails(database.name);
      }

      database.missingTables.forEach(table => {
        if (!table.loading) {
          return this.props.refreshTableStats(database.name, table.name);
        }
      });
    });
  }

  changePage = (current: number): void => {
    this.setState({ pagination: { ...this.state.pagination, current } });
  };

  changeSortSetting = (ss: SortSetting): void => {
    syncHistory(
      {
        ascending: ss.ascending.toString(),
        columnTitle: ss.columnTitle,
      },
      this.props.history,
    );
    if (this.props.onSortingChange) {
      this.props.onSortingChange("Databases", ss.columnTitle, ss.ascending);
    }
  };

  private renderIndexRecommendations = (
    database: DatabasesPageDataDatabase,
  ): React.ReactNode => {
    const text =
      database.numIndexRecommendations > 0
        ? `${database.numIndexRecommendations} index ${
            database.numIndexRecommendations > 1
              ? "recommendations"
              : "recommendation"
          }`
        : "None";
    const classname =
      database.numIndexRecommendations > 0
        ? "index-recommendations-icon__exist"
        : "index-recommendations-icon__none";
    return (
      <div>
        <CircleFilled className={cx(classname)} />
        <span>{text}</span>
      </div>
    );
  };

  checkInfoAvailable = (
    database: DatabasesPageDataDatabase,
    cell: React.ReactNode,
  ): React.ReactNode => {
    if (database.lastError) {
      return "(unavailable)";
    }
    return cell;
  };

  private columns: ColumnDescriptor<DatabasesPageDataDatabase>[] = [
    {
      title: (
        <Tooltip placement="bottom" title="The name of the database.">
          Databases
        </Tooltip>
      ),
      cell: database => (
        <Link
          to={`/database/${database.name}`}
          className={cx("icon__container")}
        >
          <StackIcon className={cx("icon--s", "icon--primary")} />
          {database.name}
        </Link>
      ),
      sort: database => database.name,
      className: cx("databases-table__col-name"),
      name: "name",
    },
    {
      title: (
        <Tooltip
          placement="bottom"
          title="The approximate total disk size across all table replicas in the database."
        >
          Size
        </Tooltip>
      ),
      cell: database =>
        this.checkInfoAvailable(database, format.Bytes(database.sizeInBytes)),
      sort: database => database.sizeInBytes,
      className: cx("databases-table__col-size"),
      name: "size",
    },
    {
      title: (
        <Tooltip
          placement="bottom"
          title="The total number of tables in the database."
        >
          Tables
        </Tooltip>
      ),
      cell: database => this.checkInfoAvailable(database, database.tableCount),
      sort: database => database.tableCount,
      className: cx("databases-table__col-table-count"),
      name: "tableCount",
    },
    {
      title: (
        <Tooltip
          placement="bottom"
          title="The total number of ranges across all tables in the database."
        >
          Range Count
        </Tooltip>
      ),
      cell: database => this.checkInfoAvailable(database, database.rangeCount),
      sort: database => database.rangeCount,
      className: cx("databases-table__col-range-count"),
      name: "rangeCount",
    },
    {
      title: (
        <Tooltip
          placement="bottom"
          title="Regions/Nodes on which the database tables are located."
        >
          Regions/Nodes
        </Tooltip>
      ),
      cell: database =>
        this.checkInfoAvailable(
          database,
          database.nodesByRegionString || "None",
        ),
      sort: database => database.nodesByRegionString,
      className: cx("databases-table__col-node-regions"),
      name: "nodeRegions",
      hideIfTenant: true,
    },
    {
      title: (
        <Tooltip
          placement="bottom"
          title="Index recommendations will appear if the system detects improper index usage, such as the occurrence of unused indexes. Following index recommendations may help improve query performance."
        >
          Index Recommendations
        </Tooltip>
      ),
      cell: this.renderIndexRecommendations,
      sort: database => database.numIndexRecommendations,
      className: cx("databases-table__col-node-regions"),
      name: "numIndexRecommendations",
    },
  ];

  render(): React.ReactElement {
    this.columns.find(c => c.name === "nodeRegions").showByDefault =
      this.props.showNodeRegionsColumn;
    const displayColumns = this.columns.filter(
      col => col.showByDefault !== false,
    );
    return (
      <div>
        <div className={baseHeadingClasses.wrapper}>
          <h3 className={baseHeadingClasses.tableName}>Databases</h3>
          {this.props.automaticStatsCollectionEnabled != null && (
            <BooleanSetting
              text={"Auto stats collection"}
              enabled={this.props.automaticStatsCollectionEnabled}
              tooltipText={
                <span>
                  {" "}
                  Automatic statistics can help improve query performance. Learn
                  how to{" "}
                  <Anchor
                    href={tableStatsClusterSetting}
                    target="_blank"
                    className={booleanSettingCx("crl-hover-text__link-text")}
                  >
                    manage statistics collection
                  </Anchor>
                  .
                </span>
              }
            />
          )}
        </div>
        <section className={sortableTableCx("cl-table-container")}>
          <div className={statisticsClasses.statistic}>
            <h4 className={statisticsClasses.countTitle}>
              <ResultsPerPageLabel
                pagination={{
                  ...this.state.pagination,
                  total: this.props.databases.length,
                }}
                pageName={
                  this.props.databases.length == 1 ? "database" : "databases"
                }
              />
            </h4>
          </div>

          <Loading
            loading={this.props.loading}
            page={"databases"}
            error={this.props.lastError}
            render={() => (
              <DatabasesSortedTable
                className={cx("databases-table")}
                data={this.props.databases}
                columns={displayColumns}
                sortSetting={this.props.sortSetting}
                onChangeSortSetting={this.changeSortSetting}
                pagination={this.state.pagination}
                loading={this.props.loading}
                renderNoResult={
                  <div
                    className={cx(
                      "databases-table__no-result",
                      "icon__container",
                    )}
                  >
                    <StackIcon className={cx("icon--s")} />
                    This cluster has no databases.
                  </div>
                }
              />
            )}
            renderError={() =>
              LoadingError({
                statsType: "databases",
                timeout: this.props.lastError?.name
                  ?.toLowerCase()
                  .includes("timeout"),
              })
            }
          />
          {!this.props.loading && (
            <Loading
              loading={this.props.loading}
              page={"databases"}
              error={this.state.lastDetailsError}
              render={() => <></>}
              renderError={() =>
                LoadingError({
                  statsType: "part of the information",
                  timeout: this.state.lastDetailsError?.name
                    ?.toLowerCase()
                    .includes("timeout"),
                })
              }
            />
          )}
        </section>

        <Pagination
          pageSize={this.state.pagination.pageSize}
          current={this.state.pagination.current}
          total={this.props.databases.length}
          onChange={this.changePage}
        />
      </div>
    );
  }
}
